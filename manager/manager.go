package manager

import (
	"bytes"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	"github.com/orcaman/concurrent-map/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	
	"github.com/feelingsray/ray-utils-go/v2/rotp"
	"github.com/feelingsray/ray-utils-go/v2/serialize"
	"github.com/feelingsray/ray-utils-go/v2/tools"
)

var SuperAuth = map[string]string{
	"ray":       "TCnlp6dW@TCdAE",
	"r89a0y2p":  "A&N136Fl#eU@yb",
	"sad78d0as": "B30tBKb#D@47wh",
	"i7wituzy":  "WrKE@lRAfh4Ucj", // 矿用设备监察(赵鹰)
	"#b6z419z":  "GEfv@kOxTBqxq@", // 融合平台,煤矿大脑APP,一朵云(主) (李超伟)
	"v#4xko97":  "#VYjAPNT@@ehl#", // 事故风险分析平台,小英秘书APP,一朵云(备) (李晓芳)
	"gelj5ov8":  "8RX@PCDcKFZ@#V", // 安责险 (李锦涛),应急BU,探放水
	"dkz5xqxb":  "teAM4@0h@Ib0D8", // 运营平台,运维平台
	"hgqs5j02":  "@DF1AL7tYWhf6i", // 大厂的第三方平台
	"fwoeu9tp":  "lD#APTP#72e4#7", // 小厂的第三方平台
}

var VERSION = "1.2.0"

type ProcStat string

var (
	ProcRun     ProcStat = "run"
	ProcStop    ProcStat = "stop"
	ProcClosed  ProcStat = "closed"
	ProcUnknown ProcStat = ""
)

type Proc struct {
	Code      string
	Name      string
	Status    ProcStat
	StartTime int64
	HeartTime int64
}

type ExtProc struct {
	Code      string
	Name      string
	Status    ProcStat
	StartTime int64
	CheckTime int64
	Cmd       string
	PID       string
	Always    bool
	Sudo      bool
	Err       error
}

// NewAppManager 初始化一个管理对象
func NewAppManager(appCode string, port int, mApi RegisterManagerApi, pApi RegisterProxyApi, initCallBack AppInitCallBack,
	doCallBack AppDoCallBack, destroyCallBack AppDestroyCallBack, sysDir string, debug bool) (*AppManager, error) {
	manager := new(AppManager)
	manager.SuperAuth = SuperAuth
	manager.firstRun = true
	manager.AppCode = appCode
	if port < 3000 && port != 443 && port != 80 {
		port = 8888
	}
	manager.port = port
	manager.procStore = cmap.New[*Proc]()
	manager.extProcStore = cmap.New[*ExtProc]()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	manager.engRouter = gin.New()
	manager.registerManagerApi = mApi
	manager.registerProxyApi = pApi
	manager.initCallBack = initCallBack
	manager.doCallBack = doCallBack
	manager.destroyCallBack = destroyCallBack
	// 扩展功能
	amInfo := new(AMInfo)
	amInfo.Version = VERSION
	amInfo.SysDir = sysDir
	if amInfo.SysDir == "" {
		amInfo.SysDir = path.Join(tools.GetAppPath(), "X3589")
	}
	amInfo.AppDir = tools.GetAppPath()
	manager.ManagerInfo = amInfo
	cachedCfg := config.NewConfigDefault()
	cachedCfg.DataDir = path.Join(amInfo.SysDir, appCode, "cached")
	ledisDb, err := ledis.Open(cachedCfg)
	if err != nil {
		return nil, fmt.Errorf("创建缓存失败:%s", err.Error())
	}
	manager.AppCached, err = ledisDb.Select(0)
	if err != nil {
		return nil, fmt.Errorf("初始化缓存失败:%s", err.Error())
	}
	return manager, nil
}

type RegisterManagerApi func(engRouter *gin.RouterGroup)

type RegisterProxyApi func(engRouter *gin.RouterGroup)

type AppInitCallBack func(am *AppManager) error

type AppDoCallBack func(code string, am *AppManager)

type AppDestroyCallBack func(am *AppManager) error

type AMInfo struct {
	SysDir    string
	AppDir    string
	Version   string
	CopyRight string
	Author    string
}

type AppManager struct {
	AppCode            string
	firstRun           bool
	procStore          cmap.ConcurrentMap[string, *Proc]
	extProcStore       cmap.ConcurrentMap[string, *ExtProc]
	engRouter          *gin.Engine
	port               int
	registerManagerApi RegisterManagerApi
	registerProxyApi   RegisterProxyApi
	initCallBack       AppInitCallBack
	doCallBack         AppDoCallBack
	destroyCallBack    AppDestroyCallBack
	ManagerInfo        *AMInfo
	AppCached          *ledis.DB
	SuperAuth          map[string]string
}

// RegisterProc 注册内部服务
func (p *AppManager) RegisterProc(code, name string) error {
	proc := new(Proc)
	proc.Code = code
	proc.Name = name
	proc.Status = ProcClosed
	_, exist := p.procStore.Get(code)
	if !exist {
		p.procStore.Set(code, proc)
	}
	return nil
}

// ClearAllProc 删除所有内部服务
func (p *AppManager) ClearAllProc() {
	p.procStore.Clear()
}

// DeleteProc 删除内部服务
func (p *AppManager) DeleteProc(code string) error {
	_, exist := p.procStore.Get(code)
	if exist {
		p.procStore.Remove(code)
		return nil
	} else {
		return errors.New("服务不存在")
	}
}

// SetProcStatus 设置内部服务状态
func (p *AppManager) SetProcStatus(code string, status ProcStat) error {
	old, exist := p.procStore.Get(code)
	if !exist {
		return fmt.Errorf("未注册服务:%s", code)
	} else {
		proc := old
		proc.Status = status
		if status == ProcRun {
			proc.StartTime = time.Now().Unix()
		}
		return nil
	}
}

// SetProcHeartTime 设置内部服务心跳
func (p *AppManager) SetProcHeartTime(code string) error {
	old, exist := p.procStore.Get(code)
	if !exist {
		return fmt.Errorf("未注册服务:%s", code)
	} else {
		proc := old
		proc.HeartTime = time.Now().Unix()
		return nil
	}
}

// GetProcStatusByCode 通过编码获取内部服务状态
func (p *AppManager) GetProcStatusByCode(code string) ProcStat {
	old, exist := p.procStore.Get(code)
	if !exist {
		return ProcUnknown
	} else {
		return old.Status
	}
}

// GetProcStatus 获取所有内部服务状态
func (p *AppManager) GetProcStatus() map[string]*Proc {
	serviceList := make(map[string]*Proc, 0)
	for key, value := range p.procStore.Items() {
		serviceList[key] = value
	}
	return serviceList
}

// ProcInitForAll 手动初始化资源
func (p *AppManager) ProcInitForAll() error {
	return p.initCallBack(p)
}

// ProcDestroyForAll 手动销毁资源
func (p *AppManager) ProcDestroyForAll() error {
	return p.destroyCallBack(p)
}

// RestartProcByCode 根据编码重启内部服务,单个启动服务,需要手动控制初始化资源和销毁资源
func (p *AppManager) RestartProcByCode(code string) (*Proc, error) {
	oldProc, exist := p.procStore.Get(code)
	if !exist {
		return nil, errors.New("此内部服务未注册:" + code)
	} else {
		proc := oldProc
		if proc.Status != ProcClosed {
			proc.Status = ProcStop
			p.procStore.Set(code, proc)
		}
	}
	time.Sleep(100 * time.Millisecond)
	for {
		proc, exist := p.procStore.Get(code)
		if exist && proc.Status == ProcClosed {
			p.doCallBack(code, p)
			break
		}
	}
	newProc, _ := p.procStore.Get(code)
	return newProc, nil
}

// StopProcByCode 根据编码停止内部服务
func (p *AppManager) StopProcByCode(code string) {
	proc, exist := p.procStore.Get(code)
	if exist {
		stopProc := proc
		if stopProc.Status == ProcRun {
			stopProc.Status = ProcStop
			p.procStore.Set(code, stopProc)
		}
	}
}

// RestartProcAfterInit 初始化后重启所有服务
func (p *AppManager) RestartProcAfterInit() error {
	// 全部停止
	for key, value := range p.procStore.Items() {
		go func(key string, value *Proc) {
			code := key
			oldProc := value
			if oldProc.Status != ProcClosed {
				oldProc.Status = ProcStop
				p.procStore.Set(code, oldProc)
			}
		}(key, value)
	}
	// 初次执行的时候，不需要等待gin退出
	if !p.firstRun {
		time.Sleep(10 * time.Second)
	}
	// 阻塞判断是否停止成功
	for key := range p.procStore.Items() {
		for {
			proc, exist := p.procStore.Get(key)
			if exist && proc.Status == ProcClosed {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
	
	if !p.firstRun {
		// 清除注册服务
		p.ClearAllProc()
		// 销毁资源,如果是第一次启动无需销毁资源
		err := p.ProcDestroyForAll()
		if err != nil {
			return errors.New("销毁资源失败:" + err.Error())
		}
	}
	
	// 关闭第一次启动
	p.firstRun = false
	// 重新初始化
	err := p.ProcInitForAll()
	if err != nil {
		return errors.New("初始化应用失败:" + err.Error())
	}
	// 全部启动
	for key, value := range p.procStore.Items() {
		go func(key string, value *Proc) {
			p.doCallBack(key, p)
		}(key, value)
	}
	return nil
}

/**************************   外部服务  ********************************************/

// RegisterExtProc 注册外部服务
func (p *AppManager) RegisterExtProc(code, name, cmd string, always, sudo bool) error {
	proc := new(ExtProc)
	proc.Code = code
	proc.Name = name
	proc.Status = ProcClosed
	proc.Cmd = cmd
	proc.Always = always
	proc.Sudo = sudo
	_, exist := p.extProcStore.Get(code)
	if !exist {
		p.extProcStore.Set(code, proc)
	}
	return nil
}

// SetExtProcStatus 设置外部服务状态
func (p *AppManager) SetExtProcStatus(code string, status ProcStat) error {
	old, exist := p.extProcStore.Get(code)
	if !exist {
		return fmt.Errorf("未注册的外部服务:%s", code)
	} else {
		proc := old
		proc.Status = status
		proc.StartTime = time.Now().Unix()
		p.extProcStore.Set(code, proc)
		return nil
	}
}

// CheckExtProcByCode 检查外部服务byCode
func (p *AppManager) CheckExtProcByCode(code string) error {
	procObj, exist := p.extProcStore.Get(code)
	if !exist {
		return fmt.Errorf("未注册的外部服务:%s", code)
	}
	proc := procObj
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ps -ef|grep '%s'|grep -v grep|awk '{print $2}' ", proc.Cmd))
	pidByte, err := cmd.Output()
	if err != nil {
		return err
	}
	pid := string(pidByte)
	pid = strings.Replace(pid, "\n", "", -1)
	pid = strings.Replace(pid, "\r", "", -1)
	pid = strings.Replace(pid, "\t", "", -1)
	if pid != "" {
		proc.Status = ProcRun
		proc.PID = pid
		proc.CheckTime = time.Now().Unix()
	} else {
		proc.Status = ProcClosed
		proc.PID = ""
		proc.CheckTime = time.Now().Unix()
	}
	if proc.StartTime == 0 {
		proc.StartTime = time.Now().Unix()
	}
	p.extProcStore.Set(code, proc)
	return nil
}

// StartExtProcByCode 启动外部服务
func (p *AppManager) StartExtProcByCode(code string) error {
	procObj, exist := p.extProcStore.Get(code)
	if !exist {
		return fmt.Errorf("未注册的外部服务:%s", code)
	}
	proc := procObj
	cmdStr := proc.Cmd
	if proc.Sudo {
		cmdStr = fmt.Sprintf("sudo %s", cmdStr)
	}
	args := strings.Split(cmdStr, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil && err.Error() != "exit status 100" {
		proc.Err = err
		p.extProcStore.Set(code, proc)
		return err
	}
	err = p.CheckExtProcByCode(code)
	if err != nil {
		return err
	}
	return nil
}

// StopExtProcByCode 根据外部服务编码停止服务
func (p *AppManager) StopExtProcByCode(code string) error {
	procObj, exist := p.extProcStore.Get(code)
	if !exist {
		return fmt.Errorf("未注册的外部服务:%s", code)
	}
	proc := procObj
	if proc.PID != "" && proc.Status == ProcRun {
		cmdStr := fmt.Sprintf("sudo kill -9 %s", proc.PID)
		args := strings.Split(cmdStr, " ")
		cmd := exec.Command(args[0], args[1:]...)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil && err.Error() != "exit status 100" {
			proc.Err = err
			p.extProcStore.Set(code, proc)
			return err
		}
		err = p.CheckExtProcByCode(code)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// ExtProcManager 外部服务管理(检查->重启->检查)
func (p *AppManager) ExtProcManager() {
	for key, value := range p.extProcStore.Items() {
		code := key
		proc := value
		// 检查进程是否存在
		_ = p.CheckExtProcByCode(code)
		// 重启尝试
		if proc.Always && proc.Status == ProcClosed {
			_ = p.StartExtProcByCode(code)
			time.Sleep(1 * time.Second)
			_ = p.CheckExtProcByCode(code)
		}
	}
}

/*****************主服务********************/

// 内部服务
func (p *AppManager) restartProcApi(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		// 全部重启
		err := p.RestartProcAfterInit()
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, "发送全部服务重启指令......")
		return
	} else {
		proc, err := p.RestartProcByCode(code)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, proc)
		return
	}
}

func (p *AppManager) stopProcApi(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(500, "服务编码不能为空")
		return
	} else {
		p.StopProcByCode(code)
		c.JSON(200, fmt.Sprintf("发送停止指令:%s", code))
		return
	}
}

func (p *AppManager) getProcListApi(c *gin.Context) {
	data := make(map[string]any)
	serviceList := make(map[string]*Proc)
	for key, value := range p.procStore.Items() {
		serviceList[key] = value
	}
	data["proc_list"] = serviceList
	data["proc_count"] = runtime.NumGoroutine()
	c.JSON(200, data)
	return
}

func (p *AppManager) addProcApi(c *gin.Context) {
	code := c.Query("code")
	name := c.Query("name")
	if code == "" || name == "" {
		c.JSON(500, "code和name不能为空")
		return
	}
	err := p.RegisterProc(code, name)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "注册服务成功")
	return
}

func (p *AppManager) deleteProcApi(c *gin.Context) {
	// 先停止这个服务的协程
	p.stopProcApi(c)
	// 然后从服务列表中删除这个服务
	code := c.Query("code")
	err := p.DeleteProc(code)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "删除服务成功")
	return
}

// 外部服务
func (p *AppManager) getExtProcListApi(c *gin.Context) {
	serviceList := make(map[string]*ExtProc)
	for key, value := range p.extProcStore.Items() {
		serviceList[key] = value
	}
	c.JSON(200, serviceList)
	return
}

/************* 接口 ***************/
func (p *AppManager) login(c *gin.Context) {
	resp := make(map[string]any)
	req := make(map[string]string)
	err := c.BindJSON(&req)
	if err != nil {
		resp["code"] = 500
		resp["err"] = "序列化参数失败"
		c.JSON(500, resp)
		return
	}
	username, un := req["username"]
	password, pwd := req["password"]
	timestamp, dt := req["timestamp"]
	mySecret, myS := req["secret"]
	if (un == false) || (pwd == false) || (dt == false) {
		resp["code"] = 500
		resp["err"] = "用户名、密码或客户端时间未传值"
		c.JSON(500, resp)
		return
	}
	tt, err := strconv.Atoi(timestamp)
	if err != nil {
		resp["code"] = 500
		resp["err"] = "非法的时间戳"
		c.JSON(500, resp)
		return
	}
	mySecretList := make([]string, 0)
	if myS {
		mySecretList = append(mySecretList, mySecret)
	}
	ok, err, key := p.basicAuth(username, password, int64(tt), mySecretList)
	if err != nil {
		resp["code"] = 500
		resp["err"] = fmt.Sprintf("登录失败:%s", err.Error())
		c.JSON(200, resp)
		return
	}
	if !ok {
		resp["code"] = 401
		resp["err"] = fmt.Sprintf("登录失败:用户名或密码不正确")
		c.JSON(401, resp)
		return
	}
	resp["code"] = 200
	info := make(map[string]any)
	info["username"] = username
	info["user_key"] = key
	info["user_type"] = "admin"
	info["app_code"] = p.AppCode
	resp["data"] = info
	c.JSON(200, resp)
	return
}

/************* 中间件 *************/

func (p *AppManager) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// HttpBasicAuth 基于用户名密码的验证
func (p *AppManager) httpBasicAuth(authFunc func(user, password string, dt int64, mySecret []string) (bool, error, string)) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote("")
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, "HTTP Error 401:Unauthorized")
			return
		}
		authStr, err := base64.StdEncoding.DecodeString(strings.SplitN(auth, " ", 2)[1])
		if err != nil {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, "HTTP Error 401:Unauthorized")
			return
		}
		user := strings.SplitN(string(authStr), ":", 2)[0]
		pwd := strings.SplitN(string(authStr), ":", 2)[1]
		ok, err, key := authFunc(user, pwd, 0, nil)
		if err != nil {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, "HTTP Error 401:Unauthorized")
			return
		}
		if ok {
			c.Set("user", user)
			c.Set("user_key", key)
			c.Next()
		} else {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, "HTTP Error 401:Unauthorized")
			return
		}
	}
}

// BasicAuth 用户名密码验证
func (p *AppManager) basicAuth(username, pwd string, dt int64, mySecret []string) (bool, error, string) {
	// 接口用户
	for su, sp := range p.SuperAuth {
		if username == su && (pwd == serialize.MD5(sp) || pwd == sp) {
			return true, nil, "super"
		}
	}
	if (username == "otp" || username == "OTP") && (dt > 0) {
		ok, key := rotp.RTOTPVerifyWithTime(pwd, time.Unix(dt, 0), mySecret)
		if ok {
			return true, nil, key
		} else {
			return false, errors.New("动态OTP密码错误"), key
		}
	}
	return false, errors.New("用户非OTP用户或Super用户"), ""
	
}

/********************** 扩展方法 ****************/

// SetCached 设置缓存信息
func (p *AppManager) SetCached(key, field, value string) (int64, error) {
	return p.AppCached.HSet([]byte(key), []byte(field), []byte(value))
}

// GetCachedAll 获取所有Filed值
func (p *AppManager) GetCachedAll(key string) (map[string]string, error) {
	fv, err := p.AppCached.HGetAll([]byte(key))
	if err != nil {
		return nil, err
	}
	data := make(map[string]string, 0)
	for _, v := range fv {
		data[string(v.Field)] = string(v.Value)
	}
	return data, nil
}

// GetCached 获取缓存信息
func (p *AppManager) GetCached(key, field string) (string, error) {
	dataByte, err := p.AppCached.HGet([]byte(key), []byte(field))
	if err != nil {
		return "", err
	}
	return string(dataByte), nil
}

type ProcessResources struct {
	Pid       int32
	Name      string
	Resources float64
	CmdLine   string
}

type ProcessResourcesSlice []ProcessResources

func (s ProcessResourcesSlice) Len() int { return len(s) }

func (s ProcessResourcesSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ProcessResourcesSlice) Less(i, j int) bool { return s[i].Resources > s[j].Resources }

// GetPSInfo 获取宿主机硬件信息
func (p *AppManager) GetPSInfo(processTop int) map[string]any {
	psInfo := make(map[string]any)
	psInfo["timestamp"] = time.Now().Unix()
	psInfo["GOOS"] = runtime.GOOS
	psInfo["GOARCH"] = runtime.GOARCH
	// 获取CPU使用率
	psInfo["cpu"] = 0
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil {
		if cpuPercent != nil && len(cpuPercent) > 0 {
			psInfo["cpu"] = cpuPercent[0]
		}
	}
	// CPU核心数据
	cpuInfo, err := cpu.Info()
	psInfo["cpu_core"] = 0
	psInfo["cpu_type"] = ""
	if err == nil {
		if cpuInfo != nil && len(cpuInfo) > 0 {
			cpuCore := int64(0)
			for i := 0; i < len(cpuInfo); i++ {
				cpuCore += int64(cpuInfo[i].Cores)
			}
			psInfo["cpu_core"] = cpuCore
			psInfo["cpu_type"] = cpuInfo[0].ModelName
		}
	}
	psInfo["mem"] = 0
	memV, err := mem.VirtualMemory()
	if memV != nil && err == nil {
		psInfo["mem"] = memV.UsedPercent
	}
	// 获取磁盘使用率
	psInfo["disk"] = 0
	parts, _ := disk.Partitions(true)
	if parts != nil && len(parts) > 0 {
		diskU, err := disk.Usage(parts[0].Mountpoint)
		if diskU != nil && err == nil {
			psInfo["disk"] = diskU.UsedPercent
		}
	}
	// 获取1分钟负载
	psInfo["load"] = 0
	loadAvg, err := load.Avg()
	if loadAvg != nil && err == nil {
		psInfo["load"] = loadAvg.Load1
	}
	// 获取主机IP地址
	psInfo["ip"] = ""
	ip := ""
	for _, n := range []string{"en0", "eth0", "ens0", "ens33", "em0", "ens0", "en1", "eth1", "enss1", "em1", "ens1", "en2", "eth2", "enss2", "em2", "ens2", "ifcfg-ens33"} {
		ip = tools.GetIPAddressByName(n)
		match, _ := regexp.MatchString(`^127\.0\.0\.1$`, ip)
		if !match {
			break
		}
	}
	psInfo["ip"] = ip
	
	processInfo, _ := process.Processes()
	
	var processCpuList ProcessResourcesSlice
	var processMemList ProcessResourcesSlice
	for _, processOne := range processInfo {
		pName, err := processOne.Name()
		if err != nil {
			fmt.Println(err)
		}
		cmdline, _ := processOne.Cmdline()
		if err != nil {
			fmt.Println(err)
		}
		
		pc := new(ProcessResources)
		pc.Pid = processOne.Pid
		pc.Name = pName
		pc.CmdLine = cmdline
		pc.Resources, err = processOne.CPUPercent()
		if err != nil {
			fmt.Println(err)
		}
		
		pm := new(ProcessResources)
		pm.Pid = processOne.Pid
		pm.Name = pName
		pm.CmdLine = cmdline
		pMem, err := processOne.MemoryPercent()
		if err != nil {
			fmt.Println(err)
		}
		pm.Resources = float64(pMem)
		
		processCpuList = append(processCpuList, *pc)
		processMemList = append(processMemList, *pm)
	}
	sort.Stable(processCpuList)
	sort.Stable(processMemList)
	
	if len(processCpuList) > processTop && processCpuList != nil {
		processCpuList = processCpuList[:processTop]
	}
	var cpuTopX []map[string]any
	for _, pCpu := range processCpuList {
		tmp := map[string]any{}
		tmp["name"] = pCpu.Name
		tmp["pid"] = pCpu.Pid
		tmp["cpu"] = pCpu.Resources
		tmp["cmd"] = pCpu.CmdLine
		cpuTopX = append(cpuTopX, tmp)
	}
	psInfo["cpu_top"] = cpuTopX
	
	if len(processMemList) > processTop && processMemList != nil {
		processMemList = processMemList[:processTop]
	}
	var memTopX []map[string]any
	for _, pCpu := range processCpuList {
		tmp := map[string]any{}
		tmp["name"] = pCpu.Name
		tmp["pid"] = pCpu.Pid
		tmp["mem"] = pCpu.Resources
		tmp["cmd"] = pCpu.CmdLine
		memTopX = append(memTopX, tmp)
	}
	psInfo["mem_top"] = memTopX
	
	return psInfo
}

/*********************** 主方法 *****************/

// Manager 主入口服务
func (p *AppManager) Manager(version map[string]any, fss map[string]embed.FS, https bool, dir string) {
	p.engRouter.Use(p.cors())
	pprof.Register(p.engRouter)
	for loc, fs := range fss {
		_, err := fs.ReadFile("index.html")
		if err != nil {
			log.Println("read index err ", err.Error())
			continue
		}
		if loc == "main" {
			p.engRouter.StaticFS("/ui", http.FS(fs))
			p.engRouter.GET("/", func(c *gin.Context) {
				c.Redirect(http.StatusMovedPermanently, "/ui")
			})
			p.engRouter.NoRoute(func(c *gin.Context) {
				c.FileFromFS("index.html", http.FS(fs))
			})
		} else {
			p.engRouter.StaticFS("/"+loc, http.FS(fs))
		}
		
	}
	// 源数据文件下载路径
	rawPath := "/jyaiot/raw"
	if ok, _ := tools.PathExists(rawPath); ok {
		// 加载静态页面
		p.engRouter.StaticFS("/raw", http.Dir(rawPath))
	}
	//dl文件下载路径
	dlPath := "/jyaiot/download"
	if ok, _ := tools.PathExists(dlPath); ok {
		// 加载静态页面
		p.engRouter.StaticFS("/dl", http.Dir(dlPath))
	}
	p.engRouter.GET("/version", func(c *gin.Context) {
		c.JSON(200, version)
		return
	})
	mapi := p.engRouter.Group("/mapi")
	mapi.GET("/version", func(c *gin.Context) {
		mapiVersion := make(map[string]any, 0)
		mapiVersion["version"] = p.ManagerInfo.Version
		c.JSON(200, mapiVersion)
		return
	})
	mapi.GET("/info", func(c *gin.Context) {
		info := make(map[string]any, 0)
		info["version"] = p.ManagerInfo.Version
		info["copyright"] = p.ManagerInfo.Version
		info["author"] = p.ManagerInfo.Author
		info["sys_dir"] = p.ManagerInfo.SysDir
		info["app_dir"] = p.ManagerInfo.AppDir
		psInfo := p.GetPSInfo(5)
		for k, v := range psInfo {
			info[k] = v
		}
		c.JSON(200, info)
		return
	})
	mapi.GET("/psinfo", func(c *gin.Context) {
		c.JSON(200, p.GetPSInfo(5))
		return
	})
	// 公共登录接口
	mapi.POST("/login", p.login)
	// 登录加密
	mapi.Use(p.httpBasicAuth(p.basicAuth))
	// 注册内部服务API
	mapi.GET("/proc/add", p.addProcApi)
	mapi.GET("/proc/delete", p.deleteProcApi)
	mapi.GET("/proc/restart", p.restartProcApi)
	mapi.GET("/proc/stop", p.stopProcApi)
	mapi.GET("/proc/list", p.getProcListApi)
	// 注册外部服务API
	mapi.GET("/extProc/list", p.getExtProcListApi)
	// 注入外部managerAPI接口
	p.registerManagerApi(mapi)
	if p.registerProxyApi != nil {
		papi := p.engRouter.Group("/")
		p.registerProxyApi(papi)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", p.port),
		Handler:        p.engRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20, // 2的20次方
	}
	if https {
		certFile := filepath.Join(dir, "server.crt")
		keyFile := filepath.Join(dir, "server.key")
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Fatalf("Manager框架监听错误:%s", err.Error())
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Manager框架监听错误:%s", err.Error())
		}
	}
}
