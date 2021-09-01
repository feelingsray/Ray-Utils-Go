package appManager

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "os/exec"
    "strings"
    "sync"
    "time"
)

/*
 * 外部服务
 */
type ExtProc struct {
    Code        string
    Name        string
    Status      ProcStat
    StartTime   time.Time
    CheckTime   time.Time
    Cmd         string
    PID         string
    Always      bool
    Sudo        bool
    Err         error
}

func NewExternalManager(port int,inv int) *ExternalManager {
    em := new(ExternalManager)
    em.procStore = new(sync.Map)
    em.port = port
    em.inv = inv
    return em
}


/*
 * 外部服务
 */
type ExternalManager struct {
    procStore       *sync.Map
    port           int
    inv             int
}

/*
 * 创建服务
 */
func (p *ExternalManager)CreateProcStore(code string,name string,cmd string,always bool,sudo bool) error {
    proc := new(ExtProc)
    proc.Code = code
    proc.Name = name
    proc.Status = ProcClosed
    proc.Cmd = cmd
    proc.Always = always
    proc.Sudo = sudo

    _,exist := p.procStore.Load(code)
    if !exist {
        p.procStore.Store(code,proc)
    }
    return nil
}

/*
 * 更新服务状态
 */
func (p *ExternalManager)UpdateProcStatus(code string,status ProcStat) error {
    old,exist := p.procStore.Load(code)
    if !exist {
        return errors.New(fmt.Sprintf("未注册的外部服务:%s",code))
    } else {
        proc := old.(*ExtProc)
        proc.Status = status
        proc.StartTime = time.Now()
        return nil
    }
}

func (p *ExternalManager)CheckerAndRestart()  {
    p.procStore.Range(func(key, value interface{}) bool {
        code := key.(string)
        proc := value.(*ExtProc)
        // 检查进程是否存在
        err := p.checkOne(proc)
        if err != nil {
            proc.Err = err
        } else {
            proc.Err = nil
        }
        // 重启尝试
        if proc.Always && proc.Status == ProcClosed{
            p.startOne(proc)
            time.Sleep(1*time.Second)
            err := p.checkOne(proc)
            if err != nil {
                proc.Err = err
            } else {
                proc.Err = nil
            }
        }
        p.procStore.Store(code,proc)
        return true
    })
}

func (p *ExternalManager)checkOne(proc *ExtProc) error {
    cmd := exec.Command("bash","-c",fmt.Sprintf("ps -ef|grep '%s'|grep -v grep|awk '{print $2}' ",proc.Cmd))
    pidByte,err := cmd.Output()
    if err != nil {
        return err
    }
    pid := string(pidByte)
    pid = strings.Replace(pid,"\n","",-1)
    pid = strings.Replace(pid,"\r","",-1)
    pid = strings.Replace(pid,"\t","",-1)
    pid = strings.Replace(pid," ","",-1)
    if pid != "" {
        proc.Status = ProcRun
        proc.PID = pid
        proc.CheckTime = time.Now()
    } else {
        proc.Status = ProcClosed
        proc.PID = ""
        proc.CheckTime = time.Now()
    }
    if proc.StartTime.IsZero()  {
        proc.StartTime = time.Now()
    }
    return nil
}

func (p *ExternalManager)startOne(proc *ExtProc) {
    cmdStr := proc.Cmd
    if proc.Sudo {
        cmdStr = fmt.Sprintf("sudo %s",cmdStr)
    }
    args := strings.Split(cmdStr," ")
    cmd := exec.Command(args[0],args[1:]...)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil && err.Error() != "exit status 100"{
        proc.Err = err
    }
}

/*
 * 停止服务
 */
func (p *ExternalManager)procList(c *gin.Context)  {
    serviceList := make(map[string]*ExtProc)
    p.procStore.Range(func(key, value interface{}) bool {
        serviceList[key.(string)]=value.(*ExtProc)
        return true
    })
    c.JSON(200,serviceList)
    return
}

/*
 * 管理接口服务
 */
func (p *ExternalManager)Manager()  {
    router := gin.New()
    router.GET("/service/", p.procList)


    go func() {
        for {
            p.CheckerAndRestart()
            time.Sleep(time.Duration(p.inv)*time.Second)
        }
    }()

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
