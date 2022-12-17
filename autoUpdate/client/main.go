package main

import (
  "fmt"
  "time"
  
  "github.com/feelingsray/Ray-Utils-Go/autoUpdate/client/src"
)

func main() {
  conf, err := src.LoadClientConf()
  if err != nil {
    fmt.Print(err.Error())
  }
  for {
    fmt.Println("# 检查更新:" + time.Now().String())
    err = src.CheckUpdate(conf, "")
    if err != nil {
      fmt.Print(err.Error())
    }
    time.Sleep(1 * time.Hour)
  }
  
}
