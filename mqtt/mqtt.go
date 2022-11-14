package mqtt

//
//import (
//	"fmt"
//	mqtt "github.com/eclipse/paho.mqtt.golang"
//	uuid "github.com/satori/go.uuid"
//	"github.com/sirupsen/logrus"
//	"strings"
//	"time"
//)
//
//type MqttHelper struct {
//	Host      string
//	Port      int
//	ClientId  string
//	CleanSession bool
//	UserName  string
//	Password  string
//	TopicList map[string]byte
//	TimeOut   int
//	KeepAlive int
//	ClientOptions *mqtt.ClientOptions
//}
//
//func NewMqttHelper(host string, username string, password string, clientId string, cleanSession bool) *MqttHelper {
//	mh := MqttHelper{}
//	mh.ClientId = clientId
//	mh.CleanSession = cleanSession
//	mh.UserName = username
//	mh.Password = password
//	mh.Host = host
//	mh.Port = 1883
//	mh.TimeOut = 3
//	mh.KeepAlive = 60
//	return &mh
//}
//
//func (m *MqttHelper) Connect(topics map[string]byte, f mqtt.MessageHandler) (mqtt.Client,error) {
//	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", m.Host, m.Port))
//	clientId := m.ClientId
//	if clientId == ""{
//		clientId = uuid.NewV1().String()
//	}
//	opts.SetClientID(clientId)
//	opts.SetAutoReconnect(true)
//	opts.SetCleanSession(m.CleanSession)
//	opts.SetUsername(m.UserName)
//	opts.SetPassword(m.Password)
//	opts.SetConnectTimeout(time.Duration(m.TimeOut) * time.Second)
//	opts.SetOnConnectHandler(m.onConnect)
//	opts.SetConnectionLostHandler(m.onConnectionLost)
//	opts.SetDefaultPublishHandler(f)
//	opts.SetKeepAlive(time.Duration(m.KeepAlive) * time.Minute)
//	m.ClientOptions = opts
//	mc := mqtt.NewClient(opts)
//	m.TopicList = topics
//	token := mc.Connect()
//	if token.Error() != nil {
//		return nil,token.Error()
//	} else {
//		return mc,nil
//	}
//}
//
//func (m *MqttHelper) ConnectForPublish() (mqtt.Client,error) {
//	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", m.Host, m.Port))
//	clientId := m.ClientId
//	if clientId == ""{
//		clientId = uuid.NewV1().String()
//	}
//	opts.SetClientID(clientId)
//	opts.SetUsername(m.UserName)
//	opts.SetPassword(m.Password)
//	m.ClientOptions = opts
//	c := mqtt.NewClient(opts)
//	if token := c.Connect(); token.Wait() && token.Error() != nil {
//		return nil,token.Error()
//	}
//	return c,nil
//}
//
//func (m *MqttHelper) Publish(client mqtt.Client, topic string, qos byte, retained bool, payload interface{}) (bool,error) {
//	token := client.Publish(topic, qos, retained, payload)
//	if token.Wait() && token.Error() != nil {
//		return false,token.Error()
//	}
//	return true,nil
//}
//
//func (m *MqttHelper) onConnect(client mqtt.Client) {
//	if token := client.SubscribeMultiple(m.TopicList, nil); token.Wait() && token.Error() != nil {
//		m.Logger.WithFields(logrus.Fields{
//			"model": "mqtt",
//		}).Error(token.Error())
//		return
//	} else {
//		tmp_topics := make([]string, len(m.TopicList))
//		i := 0
//		for k, _ := range m.TopicList {
//			fmt.Println(k)
//			tmp_topics[i] = k
//			i++
//		}
//		m.Logger.WithFields(logrus.Fields{
//			"model": "mqtt",
//		}).Infof("mqtt subscribe:%s", strings.Join(tmp_topics, ","))
//	}
//
//}
//
//func (m *MqttHelper) Disconnect(client mqtt.Client) {
//	for k, _ := range m.TopicList {
//		if token := client.Unsubscribe(k); token.Error() != nil {
//			m.Logger.WithFields(logrus.Fields{
//				"model": "mqtt",
//			}).Error(token.Error())
//		} else {
//			m.Logger.WithFields(logrus.Fields{
//				"model": "mqtt",
//			}).Infof("mqtt unsubscribe:%s", k)
//		}
//	}
//
//	client.Disconnect(250)
//	m.Logger.WithFields(logrus.Fields{
//		"model": "mqtt",
//	}).Info("mqtt disconnect")
//
//}
//
//func (m *MqttHelper) onConnectionLost(client mqtt.Client, e error) {
//	m.Logger.WithFields(logrus.Fields{
//		"model": "mqtt",
//	}).Infof("mqtt connect lost:%s", e.Error())
//}
