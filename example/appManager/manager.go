package main

import (
    "fmt"
    "github.com/feelingsray/Ray-Utils-Go/appManager"
    "github.com/gin-gonic/gin"
    "net/http"
    "runtime"
    "time"
)

// 全局配置
var GConf *Conf

func main() {
    serviceList := map[string]string {
        "api":"后端API服务",
        "PA":"程序A",
        "PB":"程序B",
        "PC":"程序C",
    }
    app,err := appManager.NewAppManager(8888,serviceList,appInitCallback,appDoCallback)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // 启动服务
    go app.Manager()

    select {

    }
}

/*
 * 初始化
 */
func appInitCallback(am *appManager.AppManager,err error)  {
    conf := NewConf()
    GConf = conf
}

/*
 * 服务启动
 */
func appDoCallback(code string, am *appManager.AppManager)  {
    switch code {
    case "api":
        go apiDemo(am)
    case "PA":
        go processDemo("PA",am)
    case "PB":
        go processDemo("PB",am)
    case "PC":
        go processDemo("PC",am)

    }
}


type Conf struct {
    Attr       time.Time
}

func NewConf() *Conf {
    conf := new(Conf)
    conf.Attr = time.Now()
    return conf
}


/*
 * 服务Demo
 */
func processDemo(code string,am *appManager.AppManager) {
    fmt.Printf("初始化信息为:%s\n",GConf.Attr)
    fmt.Println("## 启动业务:" + code)
    _ = am.UpdateProcStatus(code,appManager.ProcRun)

    defer func() {
        fmt.Println("## 停止业务:" + code)
        _ = am.UpdateProcStatus(code,appManager.ProcClosed)
        runtime.Goexit()
    }()

    for {
        if am.GetProcStatus(code) == appManager.ProcStop {
            break
        }
        // 下面是主程序
        am.SetProcHeartTime(code)
        time.Sleep(3*time.Second)
    }
    return
}

/*
 * APIDemo
 */
func apiDemo(am *appManager.AppManager)  {
    fmt.Println("## 启动业务:api")
    _ = am.UpdateProcStatus("api",appManager.ProcRun)
    router := gin.New()
    server := &http.Server{
        Addr:           fmt.Sprintf(":%d", 8000),
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   20 * time.Second,
        MaxHeaderBytes: 1 << 20, // 2的20次方
    }
    go func() {
        for {
            if am.GetProcStatus("api") == appManager.ProcStop {
                break
            }
        }
        err := server.Close()
        if err != nil {
            fmt.Println("停止api服务失败...")

        } else {
        }
    }()
    if err := server.ListenAndServe(); err != nil {
        fmt.Println(err.Error())
    }

    fmt.Println("停止api服务.....")
    _ = am.UpdateProcStatus("api",appManager.ProcClosed)
    runtime.Goexit()
}

