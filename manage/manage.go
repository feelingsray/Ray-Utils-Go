package manage

import (
	"context"
	"crypto/tls"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/expvar"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ivanlebron/ledisdb/config"
	"github.com/ivanlebron/ledisdb/ledis"
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

var VERSION = "2.0.0"

type ProcStat string

var (
	ProcRun    ProcStat = "run"
	ProcStop   ProcStat = "stop"
	ProcClosed ProcStat = "closed"
)

type Proc struct {
	Code      string
	Name      string
	Status    ProcStat
	StartTime int64
	HeartTime int64
}

// NewAppManage 初始化一个管理对象
func NewAppManage(ctx context.Context, appCode string, port int, mApi RegisterManageApi, pApi RegisterProxyApi, feApi RegisterFeApi, initCallBack AppInitCallBack,
	doCallBack AppDoCallBack, destroyCallBack AppDestroyCallBack, sysDir string, debug bool, superAuth map[string]string,
) (*AppManage, error) {
	manage := new(AppManage)
	manage.SuperAuth = superAuth
	manage.firstRun = true
	manage.AppCode = appCode
	manage.Debug = debug
	if port < 3000 && port != 443 && port != 80 {
		port = 8888
	}
	manage.port = port
	manage.procStore = cmap.New[*Proc]()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	manage.engRouter = gin.New()
	manage.registerManageApi = mApi
	manage.registerProxyApi = pApi
	manage.registerFeApi = feApi
	manage.initCallBack = initCallBack
	manage.doCallBack = doCallBack
	manage.destroyCallBack = destroyCallBack
	manage.Ctx = ctx
	// 扩展功能
	amInfo := new(AMInfo)
	amInfo.Version = VERSION
	amInfo.SysDir = sysDir
	if amInfo.SysDir == "" {
		amInfo.SysDir = filepath.Join(tools.GetAppPath(), "X3589")
		if _, err := os.Stat(amInfo.SysDir); !os.IsNotExist(err) {
			if err = os.RemoveAll(amInfo.SysDir); err != nil {
				return nil, err
			}
		}
	}
	amInfo.AppDir = tools.GetAppPath()
	manage.ManageInfo = amInfo
	cachedCfg := config.NewConfigDefault()
	cachedCfg.DataDir = filepath.Join(amInfo.SysDir, appCode, "cached")
	ledisDb, err := ledis.Open(cachedCfg)
	if err != nil {
		return nil, fmt.Errorf("创建缓存失败:%s", err.Error())
	}
	manage.AppCached, err = ledisDb.Select(0)
	if err != nil {
		return nil, fmt.Errorf("初始化缓存失败:%s", err.Error())
	}
	return manage, nil
}

type RegisterManageApi func(engRouter *gin.RouterGroup)

type RegisterProxyApi func(engRouter *gin.RouterGroup)

type RegisterFeApi func(engRouter *gin.RouterGroup)

type AppInitCallBack func(am *AppManage) error

type AppDoCallBack func(code string, am *AppManage)

type AppDestroyCallBack func(am *AppManage) error

type AMInfo struct {
	SysDir    string
	AppDir    string
	Version   string
	CopyRight string
	Author    []string
}

type AppManage struct {
	// public field
	AppCode string

	ManageInfo *AMInfo
	AppCached  *ledis.DB
	SuperAuth  map[string]string
	Debug      bool
	Ctx        context.Context

	// privacy field
	firstRun          bool
	whitelist         map[string]bool
	procStore         cmap.ConcurrentMap[string, *Proc]
	engRouter         *gin.Engine
	port              int
	registerManageApi RegisterManageApi
	registerProxyApi  RegisterProxyApi
	registerFeApi     RegisterFeApi
	initCallBack      AppInitCallBack
	doCallBack        AppDoCallBack
	destroyCallBack   AppDestroyCallBack
}

// RegisterProc 注册内部服务
func (p *AppManage) RegisterProc(code, name string) error {
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
func (p *AppManage) ClearAllProc() {
	p.procStore.Clear()
}

// SetProcStatus 设置内部服务状态
func (p *AppManage) SetProcStatus(code string, status ProcStat) {
	old, exist := p.procStore.Get(code)
	if !exist {
		return
	}
	proc := old
	proc.Status = status
	if status == ProcRun {
		proc.StartTime = time.Now().Unix()
	}
}

// SetProcHeartTime 设置内部服务心跳
func (p *AppManage) SetProcHeartTime(code string) {
	old, exist := p.procStore.Get(code)
	if !exist {
		return
	}
	proc := old
	proc.HeartTime = time.Now().Unix()
}

// ProcInitForAll 手动初始化资源
func (p *AppManage) ProcInitForAll() error {
	return p.initCallBack(p)
}

// ProcDestroyForAll 手动销毁资源
func (p *AppManage) ProcDestroyForAll() error {
	return p.destroyCallBack(p)
}

// RestartProcAfterInit 初始化后重启所有服务
func (p *AppManage) RestartProcAfterInit() error {
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

/*****************主服务********************/

func (p *AppManage) getProcListApi(c *gin.Context) {
	data := make(map[string]any)
	serviceList := make(map[string]*Proc)
	for key, value := range p.procStore.Items() {
		serviceList[key] = value
	}
	data["proc_list"] = serviceList
	data["proc_count"] = runtime.NumGoroutine()
	c.JSON(http.StatusOK, data)
	return
}

/************* 接口 ***************/
func (p *AppManage) login(c *gin.Context) {
	resp := make(map[string]any)
	req := make(map[string]string)
	if err := c.BindJSON(&req); err != nil {
		resp["code"] = http.StatusInternalServerError
		resp["msg"] = "序列化参数失败"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	username, un := req["username"]
	password, pwd := req["password"]
	timestamp, dt := req["timestamp"]
	mySecret, myS := req["secret"]
	if (!un) || (!pwd) || (!dt) {
		resp["code"] = http.StatusInternalServerError
		resp["msg"] = "用户名、密码或客户端时间未传值"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	tt, err := strconv.Atoi(timestamp)
	if err != nil {
		resp["code"] = http.StatusInternalServerError
		resp["msg"] = "非法的时间戳"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	mySecretList := make([]string, 0)
	if myS {
		mySecretList = append(mySecretList, mySecret)
	}
	ok, err, key := p.basicAuth(username, password, int64(tt), mySecretList)
	if err != nil {
		resp["code"] = http.StatusInternalServerError
		resp["msg"] = fmt.Sprintf("登录失败:%s", err.Error())
		c.JSON(http.StatusOK, resp)
		return
	}
	if !ok {
		resp["code"] = http.StatusUnauthorized
		resp["msg"] = fmt.Sprintf("登录失败:用户名或密码不正确")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	resp["code"] = http.StatusOK
	info := make(map[string]any)
	info["username"] = username
	info["user_key"] = key
	info["user_type"] = "admin"
	info["role"] = "super"
	info["app_code"] = p.AppCode
	info["ak"], info["sk"] = p.randUser()
	resp["data"] = info
	c.JSON(http.StatusOK, resp)
	return
}

func (p *AppManage) randUser() (string, string) {
	keys := make([]string, 0, len(p.SuperAuth))
	for key := range p.SuperAuth {
		keys = append(keys, key)
	}
	randomKey := keys[rand.IntN(len(keys))]
	return randomKey, p.SuperAuth[randomKey]
}

/************* 中间件 *************/

func (p *AppManage) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// HttpBasicAuth 基于用户名密码的验证
func (p *AppManage) httpBasicAuth(authFunc func(user, password string, dt int64, mySecret []string) (bool, error, string)) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote("")
	return func(c *gin.Context) {
		if _, ok := p.whitelist[c.Request.URL.Path]; ok {
			c.Next() // 在白名单内，跳过校验
			return
		}
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "no auth"})
			return
		}
		authStr, err := base64.StdEncoding.DecodeString(strings.SplitN(auth, " ", 2)[1])
		if err != nil {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "no auth"})
			return
		}
		user := strings.SplitN(string(authStr), ":", 2)[0]
		pwd := strings.SplitN(string(authStr), ":", 2)[1]
		ok, err, key := authFunc(user, pwd, 0, nil)
		if err != nil {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "no auth"})
			return
		}
		if ok {
			c.Set("user", user)
			c.Set("user_key", key)
			c.Next()
		} else {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"http": "401 no auth"})
			return
		}
	}
}

// BasicAuth 用户名密码验证
func (p *AppManage) basicAuth(username, pwd string, dt int64, mySecret []string) (bool, error, string) {
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
		}
		return false, errors.New("动态OTP密码错误"), key
	}
	return false, errors.New("用户非OTP用户或Super用户"), ""
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
func (p *AppManage) GetPSInfo(processTop int) map[string]any {
	psInfo := make(map[string]any)
	psInfo["timestamp"] = time.Now().Unix()
	psInfo["GOOS"] = runtime.GOOS
	psInfo["GOARCH"] = runtime.GOARCH
	// 获取CPU使用率
	psInfo["cpu"] = 0
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil {
		if len(cpuPercent) > 0 {
			psInfo["cpu"] = cpuPercent[0]
		}
	}
	// CPU核心数据
	cpuInfo, err := cpu.Info()
	psInfo["cpu_core"] = 0
	psInfo["cpu_type"] = ""
	if err == nil {
		if len(cpuInfo) > 0 {
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
	if len(parts) > 0 {
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

// Manage 主入口服务
func (p *AppManage) Manage(version map[string]any, fss map[string]embed.FS, https bool, dir string, whitelist map[string]bool, middleware ...gin.HandlerFunc) {
	p.engRouter.Use(p.cors())
	p.engRouter.Use(middleware...)
	if p.Debug {
		pprof.Register(p.engRouter)
		p.engRouter.GET("/debug/vars", expvar.Handler())
	}
	if whitelist != nil {
		p.whitelist = whitelist
	}
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
	// dl文件下载路径
	dlPath := "/jyaiot/download"
	if ok, _ := tools.PathExists(dlPath); ok {
		// 加载静态页面
		p.engRouter.StaticFS("/dl", http.Dir(dlPath))
	}
	p.engRouter.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version)
		return
	})
	feApi := p.engRouter.Group("/fe")
	mapi := p.engRouter.Group("/mapi")
	mapi.GET("/version", func(c *gin.Context) {
		mapiVersion := make(map[string]any)
		mapiVersion["version"] = p.ManageInfo.Version
		c.JSON(http.StatusOK, mapiVersion)
		return
	})
	mapi.GET("/info", func(c *gin.Context) {
		info := make(map[string]any)
		info["version"] = p.ManageInfo.Version
		info["copyright"] = p.ManageInfo.Version
		info["author"] = p.ManageInfo.Author
		info["sys_dir"] = p.ManageInfo.SysDir
		info["app_dir"] = p.ManageInfo.AppDir
		psInfo := p.GetPSInfo(5)
		for k, v := range psInfo {
			info[k] = v
		}
		c.JSON(http.StatusOK, info)
		return
	})
	mapi.GET("/psinfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, p.GetPSInfo(5))
		return
	})
	// 公共登录接口
	mapi.POST("/login", p.login)
	// 登录加密
	mapi.Use(p.httpBasicAuth(p.basicAuth))
	mapi.GET("/proc/list", p.getProcListApi)
	// 注入外部ManageAPI接口
	p.registerManageApi(mapi)
	// 注入外部FeAPI接口
	p.registerFeApi(feApi)
	if p.registerProxyApi != nil {
		papi := p.engRouter.Group("/")
		p.registerProxyApi(papi)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", p.port),
		Handler:        p.engRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20, // 2的20次方
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		},
	}
	if https {
		if err := server.ListenAndServeTLS(filepath.Join(dir, "server.crt"), filepath.Join(dir, "server.key")); err != nil {
			log.Fatalf("Manage框架监听错误:%s", err.Error())
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Manage框架监听错误:%s", err.Error())
		}
	}
}
