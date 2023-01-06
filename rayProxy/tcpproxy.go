package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	
	"github.com/feelingsray/Ray-Utils-Go/encode"
	"github.com/feelingsray/Ray-Utils-Go/rayProxy/services"
	"github.com/feelingsray/Ray-Utils-Go/tools"
)

type tcpArg struct {
	Name   string `yaml:"name"`
	Local  int    `yaml:"local"`
	Remote string `yaml:"remote"`
}

func main() {
	appDir := tools.GetAppPath()
	appDir = "/Users/ray/jylink/Ray-Utils-Go/rayProxy/conf"
	argList := make([]tcpArg, 0)
	_ = encode.LoadYaml(path.Join(appDir, "tcp.yaml"), &argList)
	fmt.Println(argList)
	for _, arg := range argList {
		args := services.Args{}
		args.Parent = &arg.Remote
		ss := fmt.Sprintf(":%d", arg.Local)
		args.Local = &ss
		tcpArgs := services.TCPArgs{}
		timeout := 2000
		tcpArgs.Timeout = &timeout
		parentType := "tcp"
		tcpArgs.ParentType = &parentType
		isTLS := false
		tcpArgs.IsTLS = &isTLS
		poolSize := 100
		tcpArgs.PoolSize = &poolSize
		checkInv := 3
		tcpArgs.CheckParentInterval = &checkInv
		tcpArgs.Args = args
		services.Regist(arg.Name, services.NewTCP(), tcpArgs)
	}
	runServiceList := make([]*services.ServiceItem, 0)
	for _, arg := range argList {
		service, err := services.Run(arg.Name)
		if err != nil {
			log.Fatalf("run service [%s] fail, ERR:%s", service, err)
		}
		runServiceList = append(runServiceList, service)
		
	}
	
	select {}
	
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			for _, service := range runServiceList {
				(service.S).Clean()
			}
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
