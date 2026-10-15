package base

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/danomagnum/gologix"
)

// 支持Rockwell PLC

func NewCIPClient(host string, port int, slot int, path string, emptyPath bool, timeout int, connectionSize int) *CIPClient {
	return &CIPClient{
		Host:           host,
		Port:           port,
		Slot:           slot,
		Path:           path,
		EmptyPath:      emptyPath,
		Timeout:        timeout,
		ConnectionSize: connectionSize,
	}
}

type CIPClient struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Slot           int    `json:"slot"`
	Path           string `json:"path"`
	EmptyPath      bool   `json:"empty_path"`
	Timeout        int    `json:"timeout"`
	ConnectionSize int    `json:"connection_size"`
	Client         *gologix.Client
}

func (c *CIPClient) LinerScalingSimple(value float64, scaling float64) (float64, error) {
	data := value * scaling
	return data, nil
}

func (c *CIPClient) Connect() error {
	if c.Connected() {
		return nil
	}
	client := gologix.NewClient(c.Host)
	client.AutoConnect = true
	if c.Port > 0 {
		client.Controller.Port = uint(c.Port)
	}
	if c.Timeout > 0 {
		client.SocketTimeout = time.Duration(c.Timeout) * time.Second
	}
	if c.ConnectionSize > 0 {
		client.ConnectionSize = uint16(c.ConnectionSize)
	}
	switch {
	case c.EmptyPath:
		client.Controller.Path = bytes.NewBuffer(nil)
	case c.Path != "":
		path, err := gologix.ParsePath(c.Path)
		if err != nil {
			return err
		}
		client.Controller.Path = path
	case c.Slot > 0:
		path, err := gologix.ParsePath(fmt.Sprintf("1,%d", c.Slot))
		if err != nil {
			return err
		}
		client.Controller.Path = path
	}
	if err := client.Connect(); err != nil {
		client = nil
		return err
	}
	c.Client = client
	return nil
}

func (c *CIPClient) Connected() bool {
	return c.Client != nil && c.Client.Connected()
}

func (c *CIPClient) Close() error {
	if c.Client != nil {
		err := c.Client.Disconnect()
		c.Client = nil
		return err
	}
	return nil
}

func (c *CIPClient) Read(deviceCode, tagCode string, addr string, dt string) (any, error) {
	if c.Client == nil {
		return nil, errors.New("client not connected")
	}
	switch {
	case isBoolType(dt) && isBitAccessAddr(addr):
		var value bool
		if err := c.Client.Read(addr, &value); err != nil {
			return nil, fmt.Errorf("标签(%s)读取失败:%s", addr, err.Error())
		}
		return value, nil
	case strings.ToLower(dt) == "string":
		var value string
		if err := c.Client.Read(addr, &value); err != nil {
			return nil, fmt.Errorf("标签(%s)读取失败:%s", addr, err.Error())
		}
		return value, nil
	default:
		value, err := c.Client.Read_single(addr, gologix.CIPTypeUnknown, 1)
		if err != nil {
			return nil, fmt.Errorf("标签(%s)读取失败:%s", addr, err.Error())
		}
		return value, nil
	}
}

func isBoolType(dt string) bool {
	switch strings.ToLower(dt) {
	case "bool", "bit", "boolean":
		return true
	default:
		return false
	}
}

func isBitAccessAddr(addr string) bool {
	idx := strings.LastIndex(addr, ".")
	if idx < 0 || idx == len(addr)-1 {
		return false
	}
	bit := addr[idx+1:]
	for _, r := range bit {
		if r < '0' || r > '9' {
			return false
		}
	}
	pos, err := strconv.Atoi(bit)
	return err == nil && pos >= 0 && pos <= 31
}
