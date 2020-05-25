package main

import (
	"github.com/feelingsray/Ray-Utils-Go/encode"
	"github.com/feelingsray/Ray-Utils-Go/logger"
	"github.com/feelingsray/Ray-Utils-Go/mqtt"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	log := logger.LoggerConsoleHandle(logger.DebugLevel)
	mqttH := mqtt.NewMqttHelper("192.168.111.140","real_data","password01!",log.WithFields(logrus.Fields{
		"app": "mqttPub",
	}))
	for{
		for i:=0; i<10;i++{
			mqttC := mqttH.ConnectForPublish()
			payload,_ := encode.DumpJson(map[string]interface{}{"value":i,"timestamp":time.Now().Unix()})
			mqttH.Publish(mqttC,"/raw/test/1",2,false,payload)
			mqttC.Disconnect(250)
			time.Sleep(300*time.Millisecond)
		}
	}


}
