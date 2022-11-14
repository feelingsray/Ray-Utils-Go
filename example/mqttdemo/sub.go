package main

//
//import (
//	"fmt"
//	mqtt "github.com/eclipse/paho.mqtt.golang"
//	"github.com/feelingsray/Ray-Utils-Go/logger"
//	m "github.com/feelingsray/Ray-Utils-Go/mqtt"
//	"github.com/sirupsen/logrus"
//)
//
//func main() {
//	log := logger.LoggerConsoleHandle(logger.DebugLevel)
//	mqttH := m.NewMqttHelper("192.168.111.140","real_data","password01!",log.WithFields(logrus.Fields{
//		"app": "mqttPub",
//	}))
//	topics := make(map[string]byte)
//	topics["$queue//raw/test/#"] = 2
//	mqttC := mqttH.Connect(topics, messageArrived)
//	defer mqttC.Disconnect(250)
//	select {
//
//	}
//}
//
//func messageArrived(client mqtt.Client, message mqtt.Message)  {
//	fmt.Printf("%s:%s\n",message.Topic(),string(message.Payload()))
//}
