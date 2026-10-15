package opcda_api

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huskar-t/opcda"
)

func Host(c *gin.Context) {
	ip := c.Query("ip")
	names, err := net.LookupAddr(ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器不存在"})
		return
	}
	c.JSON(http.StatusOK, names)
}

func Progs(c *gin.Context) {
	host := c.Query("host")
	var res []string
	servers, err := opcda.GetOPCServers(host)
	if err != nil {
		fmt.Println("获取服务失败", err.Error())
		return
	}
	for _, item := range servers {
		res = append(res, item.ProgID)
	}
	c.JSON(http.StatusOK, res)
}

func Tags(c *gin.Context) {
	host := c.Query("host")
	prog := c.Query("prog")
	filter := c.Query("filter")
	value, ok := GServerMap.Load(host + prog)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器不存在"})
		return
	}
	server := value.(*ServerItem)
	browser, err := server.Server.CreateBrowser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建浏览器", "details": err.Error()})
		return
	}
	defer browser.Release()
	var allTags []string
	recursiveBrowse(browser, &allTags, filter != "")
	c.JSON(http.StatusOK, allTags)
}

func InitServer(c *gin.Context) {
	host := c.Query("host")
	prog := c.Query("prog")
	if _, ok := GServerMap.Load(host + prog); !ok {
		server, err := opcda.Connect(prog, host)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "连接失败", "err": server})
			return
		}

		var tmp ServerItem
		tmp.Server = server
		tmp.Host = host
		tmp.Prog = prog
		GServerMap.Store(host+prog, &tmp)
	}
	c.JSON(http.StatusOK, "ok")
}

type InitTagRequest struct {
	Host string   `json:"host"  binding:"required"`
	Prog string   `json:"prog" binding:"required"`
	Tags []string `json:"tags" binding:"required"`
}

func InitTag(c *gin.Context) {
	var req InitTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的请求格式", "err": err.Error()})
		return
	}
	tmp, ok := GServerMap.Load(req.Host + req.Prog)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "server不存在", "err": req.Host + req.Prog})
		return
	}
	serverItem := tmp.(*ServerItem)
	groups := serverItem.Server.GetOPCGroups()
	if serverItem.Group != nil {
		_ = groups.Remove(serverItem.Group.GetServerHandle())
		serverItem.Group = nil
	}
	group, err := groups.Add(GGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "无法创建 OPC 组", "err": err.Error()})
		return
	}
	// 3. 批量添加点位
	items := group.OPCItems()
	itemList, errs, err := items.AddItems(req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加点位过程中发生系统错误", "err": err.Error()})
		return
	}
	// 4. 收集有效的服务器句柄，同时记录失败的点位
	serverHandles := make([]uint32, 0)
	errorDetails := make(map[string]string)
	for i, tag := range req.Tags {
		if errs[i] != nil {
			errorDetails[tag] = errs[i].Error()
			continue
		}
		serverHandles = append(serverHandles, itemList[i].GetServerHandle())
	}
	if len(errorDetails) != 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "发现无效点位", "err": errorDetails})
		return
	}
	serverItem.Group = group
	serverItem.Handle = serverHandles
	serverItem.Tag = req.Tags
	GServerMap.Store(req.Host+req.Prog, serverItem)
	c.JSON(http.StatusOK, "ok")
}

type TagResult struct {
	Tag       string `json:"tag"`
	Value     any    `json:"value"`
	Quality   string `json:"quality"`
	Timestamp int64  `json:"timestamp"`
}

func Read(c *gin.Context) {
	host := c.Query("host")
	prog := c.Query("prog")
	value, ok := GServerMap.Load(host + prog)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "group不存在", "err": host + prog})
		return
	}
	serverItem := value.(*ServerItem)
	status, resultErrs, err := serverItem.Group.SyncRead(opcda.OPC_DS_DEVICE, serverItem.Handle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步读取失败: " + err.Error()})
		return
	}
	res := make([]TagResult, 0)
	for i, item := range serverItem.Tag {
		var resOne TagResult
		resOne.Tag = item
		resOne.Timestamp = status[i].Timestamp.Unix()
		if resultErrs[i] != nil {
			resOne.Quality = "bad"
			resOne.Value = "-9999"
		} else {
			resOne.Quality = "good"
			resOne.Value = status[i].Value
		}
		res = append(res, resOne)
	}
	c.JSON(http.StatusOK, res)
}

func Write(c *gin.Context) {
	host := c.Query("host")
	prog := c.Query("prog")
	tag := c.Query("tag")
	val := c.Query("value")
	value, ok := GServerMap.Load(host + prog)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器不存在"})
		return
	}
	server := value.(*ServerItem)

	// 2. 创建临时组 (读取/写入在 OPC DA 中都建议通过 Group 进行)
	groups := server.Server.GetOPCGroups()
	groupName := fmt.Sprintf("SingleWrite_%d", time.Now().UnixNano())
	group, err := groups.Add(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建 OPC 组失败: " + err.Error()})
		return
	}
	defer groups.Remove(group.GetServerHandle())

	// 3. 将点位添加到组
	itemsObj := group.OPCItems()
	itemList, errs, err := itemsObj.AddItems([]string{tag})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加点位系统错误: " + err.Error()})
		return
	}

	// 检查点位是否添加成功 (比如路径是否存在)
	if errs[0] != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"tag":   tag,
			"error": "点位无效或不存在: " + errs[0].Error(),
		})
		return
	}

	// 4. 执行写入
	// 直接调用 item 对象的 Write 方法 (对应你 demo 里的 item.Write)
	err = itemList[0].Write(val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"tag":   tag,
			"error": "写入数据失败: " + err.Error(),
		})
		return
	}

	// 5. 返回成功
	c.JSON(http.StatusOK, gin.H{
		"tag":     tag,
		"message": "写入成功",
		"value":   val,
	})
}

func Heart(c *gin.Context) {
	host := c.Query("host")
	prog := c.Query("prog")
	value, ok := GServerMap.Load(host + prog)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器不存在"})
		return
	}
	serverItem := value.(*ServerItem)
	state, err := serverItem.Server.GetServerState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取状态错误: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, state)
}
