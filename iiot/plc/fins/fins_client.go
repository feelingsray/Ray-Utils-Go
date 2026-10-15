package fins

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/iwanlebron/gofins/fins"
)

// Omron PLC NJ501-1300.
func NewFINSClient(host string, port int, localHost string, localPort int, network int, node int, unit int, localNetwork int, localNode int, localUnit int, timeout int, byteOrder string) *FINSClient {
	return &FINSClient{
		Host:         host,
		Port:         port,
		LocalHost:    localHost,
		LocalPort:    localPort,
		Network:      network,
		Node:         node,
		Unit:         unit,
		LocalNetwork: localNetwork,
		LocalNode:    localNode,
		LocalUnit:    localUnit,
		Timeout:      timeout,
		ByteOrder:    byteOrder,
		order:        binary.BigEndian,
	}
}

type FINSClient struct {
	Host         string       `json:"host"`
	Port         int          `json:"port"`
	LocalHost    string       `json:"local_host"`
	LocalPort    int          `json:"local_port"`
	Network      int          `json:"network"`
	Node         int          `json:"node"`
	Unit         int          `json:"unit"`
	LocalNetwork int          `json:"local_network"`
	LocalNode    int          `json:"local_node"`
	LocalUnit    int          `json:"local_unit"`
	Timeout      int          `json:"timeout"`
	ByteOrder    string       `json:"byte_order"`
	Client       *fins.Client `json:"-"`
	order        binary.ByteOrder
}

func (f *FINSClient) LinerScalingSimple(value float64, scaling float64) (float64, error) {
	data := value * scaling
	return data, nil
}

func (f *FINSClient) Connect() error {
	localHost := f.LocalHost
	if localHost == "" {
		localHost = "0.0.0.0"
	}
	port := f.Port
	if port == 0 {
		port = 9600
	}
	localAddr := fins.NewAddress(localHost, f.LocalPort, byte(f.LocalNetwork), byte(f.LocalNode), byte(f.LocalUnit))
	plcAddr := fins.NewAddress(f.Host, port, byte(f.Network), byte(f.Node), byte(f.Unit))
	client, err := fins.NewClient(localAddr, plcAddr)
	if err != nil {
		return err
	}
	if f.Timeout > 0 {
		client.SetTimeoutMs(uint(f.Timeout * 1000))
	}
	if strings.ToLower(f.ByteOrder) == "little" {
		client.SetByteOrder(binary.LittleEndian)
		f.order = binary.LittleEndian
	} else {
		client.SetByteOrder(binary.BigEndian)
		f.order = binary.BigEndian
	}
	f.Client = client
	return nil
}

func (f *FINSClient) Close() error {
	if f.Client != nil {
		f.Client.Close()
		f.Client = nil
	}
	return nil
}

func (f *FINSClient) Read(deviceCode, tagCode string, memoryArea string, startAddr int, bitOffset int, readLength int, dt string) (any, error) {
	area, err := f.memoryArea(memoryArea)
	if err != nil {
		return nil, fmt.Errorf("device %s, tag %s, %w", deviceCode, tagCode, err)
	}
	switch strings.ToLower(dt) {
	case "bool", "bit":
		if readLength <= 0 {
			readLength = 1
		}
		value, err := f.Client.ReadBits(area, uint16(startAddr), byte(bitOffset), uint16(readLength))
		if err != nil {
			return nil, fmt.Errorf("内存区(%s-%d.%d-%d)读取失败:%s", memoryArea, startAddr, bitOffset, readLength, err.Error())
		}
		if readLength == 1 && len(value) > 0 {
			return value[0], nil
		}
		return value, nil
	case "string":
		if readLength <= 0 {
			return nil, fmt.Errorf("内存区(%s-%d)读取长度错误", memoryArea, startAddr)
		}
		value, err := f.Client.ReadString(area, uint16(startAddr), uint16(readLength))
		if err != nil {
			return nil, fmt.Errorf("内存区(%s-%d-%d)读取失败:%s", memoryArea, startAddr, readLength, err.Error())
		}
		return value, nil
	case "byte":
		raw, err := f.readBytes(area, memoryArea, startAddr, 1)
		if err != nil {
			return nil, err
		}
		if len(raw) == 0 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return raw[0], nil
	case "bytes":
		if readLength <= 0 {
			return nil, fmt.Errorf("内存区(%s-%d)读取长度错误", memoryArea, startAddr)
		}
		return f.readBytes(area, memoryArea, startAddr, readLength)
	case "word", "uint16":
		value, err := f.Client.ReadWords(area, uint16(startAddr), 1)
		if err != nil {
			return nil, fmt.Errorf("内存区(%s-%d)读取失败:%s", memoryArea, startAddr, err.Error())
		}
		if len(value) != 1 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return value[0], nil
	case "int", "int16":
		value, err := f.Client.ReadWords(area, uint16(startAddr), 1)
		if err != nil {
			return nil, fmt.Errorf("内存区(%s-%d)读取失败:%s", memoryArea, startAddr, err.Error())
		}
		if len(value) != 1 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return int16(value[0]), nil
	case "dword", "uint32", "udint":
		raw, err := f.readBytes(area, memoryArea, startAddr, 2)
		if err != nil {
			return nil, err
		}
		if len(raw) != 4 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return f.order.Uint32(raw), nil
	case "dint", "int32":
		raw, err := f.readBytes(area, memoryArea, startAddr, 2)
		if err != nil {
			return nil, err
		}
		if len(raw) != 4 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return int32(f.order.Uint32(raw)), nil
	case "real", "float", "float32":
		raw, err := f.readBytes(area, memoryArea, startAddr, 2)
		if err != nil {
			return nil, err
		}
		if len(raw) != 4 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return math.Float32frombits(f.order.Uint32(raw)), nil
	case "lword", "uint64", "ulint":
		raw, err := f.readBytes(area, memoryArea, startAddr, 4)
		if err != nil {
			return nil, err
		}
		if len(raw) != 8 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return f.order.Uint64(raw), nil
	case "lint", "int64":
		raw, err := f.readBytes(area, memoryArea, startAddr, 4)
		if err != nil {
			return nil, err
		}
		if len(raw) != 8 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return int64(f.order.Uint64(raw)), nil
	case "lreal", "double", "float64":
		raw, err := f.readBytes(area, memoryArea, startAddr, 4)
		if err != nil {
			return nil, err
		}
		if len(raw) != 8 {
			return nil, fmt.Errorf("内存区(%s-%d)返回数据量错误", memoryArea, startAddr)
		}
		return math.Float64frombits(f.order.Uint64(raw)), nil
	default:
		return nil, fmt.Errorf("device %s, tag %s, 未知数据类型:%s", deviceCode, tagCode, dt)
	}
}

func (f *FINSClient) readBytes(area byte, memoryArea string, startAddr int, readLength int) ([]byte, error) {
	if readLength <= 0 {
		return nil, fmt.Errorf("内存区(%s-%d)读取长度错误", memoryArea, startAddr)
	}
	raw, err := f.Client.ReadBytes(area, uint16(startAddr), uint16(readLength))
	if err != nil {
		return nil, fmt.Errorf("内存区(%s-%d-%d)读取失败:%s", memoryArea, startAddr, readLength, err.Error())
	}
	return raw, nil
}

func (f *FINSClient) memoryArea(memoryArea string) (byte, error) {
	switch strings.ToLower(memoryArea) {
	case "wr_bit":
		return fins.MemoryAreaWRBit, nil
	case "wr_word":
		return fins.MemoryAreaWRWord, nil
	case "hr_bit":
		return fins.MemoryAreaHRBit, nil
	case "hr_word":
		return fins.MemoryAreaHRWord, nil
	case "ar_bit":
		return fins.MemoryAreaARBit, nil
	case "ar_word":
		return fins.MemoryAreaARWord, nil
	case "dm_bit":
		return fins.MemoryAreaDMBit, nil
	case "dm_word":
		return fins.MemoryAreaDMWord, nil
	case "cio_bit", "cio_word", "timer_counter_completion_flag", "timer_counter_pv", "task_bit", "task_status", "index_register_pv", "data_register_pv", "clock_pulses_condition_flags_bit":
		return 0, fmt.Errorf("gofins 当前读接口不支持内存区:%s", memoryArea)
	default:
		return 0, fmt.Errorf("未知内存区:%s", memoryArea)
	}
}
