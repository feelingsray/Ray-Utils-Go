//go:build windows

package opcda_api

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ole/go-ole" // 假设你使用的是这个库
	"github.com/huskar-t/opcda"
	"golang.org/x/sys/windows/svc"
)

var runPort = 11600

type myHandler struct{}

func (m *myHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	// 1. 在这里执行你原本 main 里的初始化逻辑
	runtime.LockOSThread()
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		fmt.Printf("COM 初始化状态: %v\n", err)
	}
	// 注意：这里的 defer ole.CoUninitialize() 会在服务退出时执行

	go ReConnect()
	// 这里的 defer Close() 建议根据你的业务逻辑放在合适位置

	// 2. 启动 Gin (必须开启协程，否则会阻塞 Execute 导致服务启动超时)
	go func() {
		gin.SetMode(gin.ReleaseMode) // 服务模式建议关闭 debug 日志
		router := gin.Default()
		registerRoutes(router)
		_ = router.Run(fmt.Sprintf(":%d", runPort))
	}()

	// 3. 告诉 Windows 服务已就绪
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	// 4. 监听退出信号
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				// 这里可以执行清理操作，例如 Close()
				Close()
				ole.CoUninitialize()
				return
			default:
			}
		}
	}
}

func OPCDAAPIServer(runPort int) {
	// 判断运行环境
	inService, err := svc.IsAnInteractiveSession()
	if err != nil {
		return
	}

	if !inService {
		// 以服务模式运行
		// 注意："MyService" 必须和 sc create 时的名字完全匹配
		if err := svc.Run("MyService", &myHandler{}); err != nil {
			fmt.Printf("服务运行失败: %v\n", err)
		}
		return
	}

	// --- 以下是本地直接调试运行的逻辑 ---
	runtime.LockOSThread()
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		fmt.Printf("COM 初始化状态: %v\n", err)
	}
	defer ole.CoUninitialize()

	go ReConnect()
	defer Close()

	r := gin.Default()
	// ... 注册路由 (建议把注册路由抽离成一个公共函数)
	registerRoutes(r)
	_ = r.Run(fmt.Sprintf(":%d", runPort))
}

// 建议抽离路由，方便 main 和 Execute 复用
func registerRoutes(r *gin.Engine) {
	r.GET("/host", Host)
	r.GET("/prog", Progs)
	r.GET("/tags", Tags)
	r.GET("/initServer", InitServer)
	r.POST("/initTag", InitTag)
	r.GET("/read", Read)
	r.GET("/write", Write)
	r.GET("/heart", Heart)
}

func Close() {
	GServerMap.Range(func(k, v any) bool {
		server := v.(*ServerItem)
		if server.Server == nil {
			return true
		}
		groups := server.Server.GetOPCGroups()
		if groups == nil {
			return true
		}
		_ = groups.Remove(server.Group.GetServerHandle())
		_ = server.Server.Disconnect()
		return true
	})
}

func ReConnect() {
	for {
		GServerMap.Range(func(k, v any) bool {
			server := v.(*ServerItem)
			if server.Server == nil {
				return true
			}
			state, err := server.Server.GetServerState()
			if err == nil && state == 1 {
				return true
			}
			_ = server.Server.Disconnect()
			newServer, err := opcda.Connect(server.Host, server.Prog)
			if err != nil {
				fmt.Printf("链接失败: %v\n", err)
				return true
			}
			GServerMap.Store(server.Host+server.Prog, newServer)
			return true
		})
		time.Sleep(60 * time.Second)
	}
}
