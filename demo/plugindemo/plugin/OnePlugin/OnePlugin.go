package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type OnePlugin struct {

}

func (p *OnePlugin) Worker(data map[string]interface{})  {
	redisC := data["redis"].(*redis.Pool).Get()
	_, _ = redisC.Do("SELECT", 3)
	replay, err := redis.String(redisC.Do("GET", "real:/up/tag/1/11/AYeQQJeAR/140000018"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(replay)
	fmt.Println("这是第一个")
}

func Plugin(data map[string]interface{})  {
	p := OnePlugin{}
	p.Worker(data)
}