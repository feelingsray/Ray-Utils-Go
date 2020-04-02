package main

import (
	"fmt"
	"github.com/feelingsray/Ray-Utils-Go/autoupdate/edge/src"
	"github.com/feelingsray/Ray-Utils-Go/tools"
	"net/http"
	"path"
	"time"
)

func main() {
	conf,err := src.LoadEdgeConf()
	if err != nil {
		fmt.Print(err.Error())
	}

	go func() {
		for {
			fmt.Println("# 检查更新:"+time.Now().String())
			src.EdgeDownload(conf)
			if err != nil {
				fmt.Print(err.Error())
			}
			time.Sleep(1*time.Hour)
		}
	}()

	releasePath := path.Join(tools.GetAppPath(),"releases")
	releasePath = "/Users/ray/jylink/Ray-Utils-Go/autoupdate/edge/releases"
	http.Handle("/", http.FileServer(http.Dir(releasePath)))
	http.ListenAndServe(":9091", nil)

}
