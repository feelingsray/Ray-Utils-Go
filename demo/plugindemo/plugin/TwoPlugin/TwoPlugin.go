package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type TwoPlugin struct {
}

func (p *TwoPlugin) Worker(data map[string]interface{}) {
	redisC := data["redis"].(*redis.Pool).Get()
	_, _ = redisC.Do("SELECT", 3)
	replay, _ := redis.String(redisC.Do("GET", "real:/up/tag/1/11/AYeQQJeAR/140000189"))
	fmt.Println(replay)
	fmt.Println("这是第二个")
}

func Plugin(data map[string]interface{}) {
	p := TwoPlugin{}
	p.Worker(data)
}
