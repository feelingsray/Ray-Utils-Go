package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"plugin"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/feelingsray/Ray-Utils-Go/tools"
)

func main() {
	appDir := tools.GetAppPath()
	appDir = "/Users/ray/jylink/Ray-Utils-Go/demo/plugindemo"
	pluginDir := path.Join(appDir, "pkg")
	dirList, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	args := make(map[string]interface{})

	// redis连接池
	redisPool := &redis.Pool{
		MaxIdle:     1000, /*最大的空闲连接数*/
		MaxActive:   1000, /*最大的激活连接数*/
		IdleTimeout: time.Duration(3) * time.Second,
		Dial: func() (redis.Conn, error) {
			// 访问本地,端口固定
			redisAddr := "192.168.111.140:6379"
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			if _, authErr := c.Do("AUTH", "password01!"); authErr != nil {
				c.Close()
				return nil, authErr
			}
			return c, nil
		},
	}
	defer redisPool.Close()

	args["redis"] = redisPool

	for {
		for i, pluginFile := range dirList {
			pluginPath := path.Join(pluginDir, pluginFile.Name())
			so, err := plugin.Open(pluginPath)
			if err != nil {
				log.Fatal(err.Error())
			}
			p, _ := so.Lookup("Plugin")
			f := p.(func(map[string]interface{}))
			args["plugin"] = fmt.Sprintf("%d", i)
			f(args)
			time.Sleep(1 * time.Second)
		}
	}

}
