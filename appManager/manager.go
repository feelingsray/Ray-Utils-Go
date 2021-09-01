package appManager

import (
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "sync"
    "time"
)


type ProcStat string

var (
    ProcRun ProcStat = "run"
    ProcStop ProcStat = "stop"
    ProcClosed ProcStat = "closed"
    ProcUnknown ProcStat = ""
)


/*
 * 内部服务
 */
type Proc struct {
    Code    string
    Name    string
    Status  ProcStat
    StartTime time.Time
    HeartTime time.Time
}

/*
 * 创建管理对象
 * port  端口号
 * procList  服务列表
 * initCallBack  初始化回调函数
 * doCallBack    服务执行回调函数
 */
func NewAppManager(port int,procList map[string]string,initCallBack AppInitCallBack,doCallBack AppDoCallBack) (*AppManager,error) {
    app := new(AppManager)
    app.ProcStore = new(sync.Map)
    app.port = port
    app.initCallBack = initCallBack
    app.doCallBack = doCallBack
    for code,name := range procList {
        err := app.CreateProcStore(code, name)
        if err != nil {
            return nil,err
        }
    }
    return app,nil
}

type AppDoCallBack func(code string,am *AppManager)
type AppInitCallBack func(am *AppManager,err error)

type AppManager struct {
    ProcStore       *sync.Map
    port           int
    initCallBack      AppInitCallBack
    doCallBack        AppDoCallBack

}

/*
 * 创建服务
 */
func (p *AppManager)CreateProcStore(code string,name string) error {
    proc := new(Proc)
    proc.Code = code
    proc.Name = name
    proc.Status = ProcClosed
    _,exist := p.ProcStore.Load(code)
    if !exist {
        p.ProcStore.Store(code,proc)
    }
    return nil
}

/*
 * 更新服务状态
 */
func (p *AppManager)UpdateProcStatus(code string,status ProcStat) error {
    old,exist := p.ProcStore.Load(code)
    if !exist {
        return errors.New(fmt.Sprintf("未注册服务:%s",code))
    } else {
        proc := old.(*Proc)
        proc.Status = status
        proc.StartTime = time.Now()
        return nil
    }
}

/*
 * 更新程序状态
 */
func (p *AppManager)SetProcHeartTime(code string) error {
    old,exist := p.ProcStore.Load(code)
    if !exist {
        return errors.New(fmt.Sprintf("未注册服务:%s",code))
    } else {
        proc := old.(*Proc)
        proc.HeartTime = time.Now()
        return nil
    }
}

/*
 * 获取当前状态
 */
func (p *AppManager)GetProcStatus(code string) ProcStat {
    old,exist := p.ProcStore.Load(code)
    if !exist {
        return ProcUnknown
    } else {
        return old.(*Proc).Status
    }
}

/*
 * 重启服务
 */
func (p *AppManager)ProcRestart(code string,name string) *Proc {
    oldProc,exist := p.ProcStore.Load(code)
    if !exist {
        p.CreateProcStore(code,name)
    } else {
        proc := oldProc.(*Proc)
        if proc.Status != ProcClosed {
            proc.Status = ProcStop
            p.ProcStore.Store(code,proc)
        }
    }
    time.Sleep(1*time.Second)
    for {
        proc,exist := p.ProcStore.Load(code)
        if exist && proc.(*Proc).Status == ProcClosed {
            p.doCallBack(code,p)
            break
        }
    }
    newProc,_ := p.ProcStore.Load(code)
    return newProc.(*Proc)
}

/*
 * 全部重启(带全局初始化)
 */
func (p *AppManager)AllProcRestart() (map[string]*Proc,error) {
    // 初始化
    err := *new(error)
    p.initCallBack(p,err)
    if err != nil {
        return nil,err
    }
    // 全部重启
    newProcList := make(map[string]*Proc,0)
    p.ProcStore.Range(func(key, value interface{}) bool {
        code := key.(string)
        oldProc := value.(*Proc)
        if oldProc.Status != ProcClosed {
            oldProc.Status = ProcStop
            p.ProcStore.Store(code,oldProc)
        }
        time.Sleep(1*time.Second)
        for {
            proc,exist := p.ProcStore.Load(code)
            if exist && proc.(*Proc).Status == ProcClosed {
                p.doCallBack(code,p)
                break
            }
        }
        newProc,_ := p.ProcStore.Load(code)
        newProcList[code] = newProc.(*Proc)
        return true
    })
    return newProcList,nil
}


/*
 * 停止服务
 */
func (p *AppManager)ProcStop(code string) {
    proc,exist := p.ProcStore.Load(code)
    if exist {
        stopProc := proc.(*Proc)
        if stopProc.Status == ProcRun {
            stopProc.Status = ProcStop
            p.ProcStore.Store(code,stopProc)
        }
    }
}

/*
 * 服务状态清单
 */
func (p *AppManager)ProcStatusList() map[string]*Proc{
    serviceList := make(map[string]*Proc,0)
    p.ProcStore.Range(func(key, value interface{}) bool {
        serviceList[key.(string)]=value.(*Proc)
        return true
    })
    return serviceList
}

/*
 * 管理接口服务
 */
func (p *AppManager)Manager()  {
    router := gin.New()
    router.GET("/restart/",p.procRestart)
    router.GET("/stop/", p.procStop)
    router.GET("/service/", p.procList)

    server := &http.Server{
        Addr:           fmt.Sprintf(":%d", p.port),
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   20 * time.Second,
        MaxHeaderBytes: 1 << 20, // 2的20次方
    }
    if err := server.ListenAndServe(); err != nil {
        fmt.Println(err.Error())
    }
}


