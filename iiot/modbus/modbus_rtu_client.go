package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/ivanlebron/modbus"
)

// NewModbusRTUClient 初始化 RTU 客户端
func NewModbusRTUClient(address string, baudRate int, slaverId int, timeout int) *ModbusRTUClient {
	mc := new(ModbusRTUClient)
	mc.Address = address
	mc.BaudRate = baudRate
	mc.SlaveId = slaverId
	mc.Timeout = timeout
	// 默认 8N1，如果需要灵活配置，可以在此处增加参数
	mc.DataBits = 8
	mc.StopBits = 1
	mc.Parity = "N"
	return mc
}

type ModbusRTUClient struct {
	Address    string `json:"address"` // 串口地址，例如 COM1 或 /dev/ttyUSB0
	BaudRate   int    `json:"baud_rate"`
	DataBits   int    `json:"data_bits"`
	StopBits   int    `json:"stop_bits"`
	Parity     string `json:"parity"`
	SlaveId    int    `json:"slave_id"`
	Timeout    int    `json:"timeout"`
	RTUHandler *modbus.RTUClientHandler
	Client     modbus.Client
}

func (m *ModbusRTUClient) LinerScalingSimple(value float64, scaling float64) (float64, error) {
	data := value * scaling
	return data, nil
}

func (m *ModbusRTUClient) LinerScaling(value any, offset float64, rawHeight float64, rawLow float64, scaleHeight float64, scaleLow float64) (float64, error) {
	valueType := reflect.TypeOf(value).Name()
	switch valueType {
	case "string":
		return 0, errors.New("can not a string param")
	case "bool":
		return 0, errors.New("can not a bool param")
	}
	// 这里还没做量程判断
	data := (value.(float64)-offset)/(rawHeight-rawLow)*(scaleHeight-scaleLow) + scaleLow
	return data, nil
}

func (m *ModbusRTUClient) ByteTo(deviceCode, tagCode string, data []byte, dataType string, displayType string) (any, error) {
	if dataType == "float" {
		if len(data) != 4 {
			return nil, errors.New("数据长度不正确,应为4个8位")
		}
		if displayType == "ABCD" {
			bits := binary.BigEndian.Uint32(data)
			return math.Float32frombits(bits), nil
		}
		if displayType == "CDAB" {
			bits := binary.BigEndian.Uint32([]byte{data[2], data[3], data[0], data[1]})
			return math.Float32frombits(bits), nil
		}
		if displayType == "DCBA" {
			bits := binary.BigEndian.Uint32([]byte{data[3], data[2], data[1], data[0]})
			return math.Float32frombits(bits), nil
		}
		if displayType == "BADC" {
			bits := binary.BigEndian.Uint32([]byte{data[1], data[0], data[3], data[1]})
			return math.Float32frombits(bits), nil
		}
	} else if dataType == "long" {
		if len(data) != 4 {
			return nil, errors.New("数据长度不正确,应为4个8位")
		}
		if displayType == "ABCD" {
			bits := binary.BigEndian.Uint32(data)
			return int32(bits), nil
		}
		if displayType == "CDAB" {
			bits := binary.BigEndian.Uint32([]byte{data[2], data[3], data[0], data[1]})
			return int32(bits), nil
		}
		if displayType == "DCBA" {
			bits := binary.BigEndian.Uint32([]byte{data[3], data[2], data[1], data[0]})
			return int32(bits), nil
		}
		if displayType == "BADC" {
			bits := binary.BigEndian.Uint32([]byte{data[1], data[0], data[3], data[2]})
			return int32(bits), nil
		}
	} else if dataType == "double" {
		if len(data) != 8 {
			return nil, errors.New("数据长度不正确,应为8个8位")
		}
		if displayType == "ABCDEFGH" {
			bits := binary.BigEndian.Uint64(data)
			return math.Float64frombits(bits), nil
		}
		if displayType == "GHEFCDAB" {
			bits := binary.BigEndian.Uint64([]byte{data[6], data[7], data[4], data[5], data[2], data[3], data[0], data[1]})
			return math.Float64frombits(bits), nil
		}
		if displayType == "BADCFEHG" {
			bits := binary.BigEndian.Uint64([]byte{data[1], data[0], data[3], data[2], data[5], data[4], data[7], data[6]})
			return math.Float64frombits(bits), nil
		}
		if displayType == "HGFEDCBA" {
			bits := binary.BigEndian.Uint64([]byte{data[7], data[6], data[5], data[4], data[3], data[2], data[1], data[0]})
			return math.Float64frombits(bits), nil
		}
	} else {
		if len(data) != 2 {
			return nil, errors.New("数据长度不正确,应为2个8位")
		}
		if displayType == "AB" {
			bits := binary.BigEndian.Uint16(data)
			return int16(bits), nil
		}
		bits := binary.LittleEndian.Uint16(data)
		return int16(bits), nil
	}
	return nil, nil
}

func (m *ModbusRTUClient) ToByte(data any, dataType string, displayType string) ([]byte, error) {
	if dataType == "float" {
		if reflect.TypeOf(data).Kind() != reflect.Float32 {
			return nil, errors.New("数据类型不匹配")
		}
		if displayType == "ABCD" {
			bits := math.Float32bits(data.(float32))
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, bits)
			return bytes, nil
		}
		if displayType == "CDAB" {
			bits := math.Float32bits(data.(float32))
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, bits)
			return []byte{bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		}
		if displayType == "DCBA" {
			bits := math.Float32bits(data.(float32))
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, bits)
			return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		}
		if displayType == "BADC" {
			bits := math.Float32bits(data.(float32))
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, bits)
			return []byte{bytes[1], bytes[0], bytes[3], bytes[1]}, nil
		}
	} else if dataType == "long" {
		if reflect.TypeOf(data).Kind() != reflect.Int32 {
			return nil, errors.New("数据类型不匹配")
		}
		if displayType == "ABCD" {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uint32(data.(int32)))
			return bytes, nil
		}
		if displayType == "CDAB" {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uint32(data.(int32)))
			return []byte{bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		}
		if displayType == "DCBA" {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uint32(data.(int32)))
			return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		}
		if displayType == "BADC" {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uint32(data.(int32)))
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2]}, nil
		}
	} else if dataType == "double" {
		if reflect.TypeOf(data).Kind() != reflect.Float64 {
			return nil, errors.New("数据类型不匹配")
		}
		if displayType == "ABCDEFGH" {
			bits := math.Float64bits(data.(float64))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return bytes, nil
		}
		if displayType == "GHEFCDAB" {
			bits := math.Float64bits(data.(float64))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return []byte{bytes[6], bytes[7], bytes[4], bytes[5], bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		}
		if displayType == "BADCFEHG" {
			bits := math.Float64bits(data.(float64))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2], bytes[5], bytes[4], bytes[7], bytes[6]}, nil
		}
		if displayType == "HGFEDCBA" {
			bits := math.Float64bits(data.(float64))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return []byte{bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		}
	} else {
		if reflect.TypeOf(data).Kind() != reflect.Int16 {
			return nil, errors.New("数据类型不匹配")
		}
		if displayType == "AB" {
			bytes := make([]byte, 2)
			binary.BigEndian.PutUint16(bytes, uint16(data.(int16)))
			return bytes, nil
		}
		bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytes, uint16(data.(int16)))
		return bytes, nil
	}
	return nil, nil
}

func (m *ModbusRTUClient) Connect() error {
	handler := modbus.NewRTUClientHandler(m.Address)
	handler.BaudRate = m.BaudRate
	handler.DataBits = m.DataBits
	handler.Parity = m.Parity
	handler.StopBits = m.StopBits
	handler.SlaveId = byte(m.SlaveId)
	handler.Timeout = time.Duration(m.Timeout) * time.Second
	handler.IdleTimeout = time.Duration(m.Timeout) * time.Second

	err := handler.Connect()
	if err != nil {
		m.Client = nil
		return err
	}
	m.RTUHandler = handler
	m.Client = modbus.NewClient(handler)
	return nil
}

// Close 关闭串口
func (m *ModbusRTUClient) Close() error {
	if m.Client != nil {
		m.Client = nil
	}
	if m.RTUHandler != nil {
		return m.RTUHandler.Close()
	}
	return nil
}

// Reconnect 检查并重新连接
func (m *ModbusRTUClient) Reconnect() error {
	// 1. 如果 Client 或 Handler 为空，说明从未连接，直接尝试连接
	if m.Client == nil || m.RTUHandler == nil {
		return m.Connect()
	}
	// 2. 尝试进行一次极轻量的读取操作（读取线圈 0），测试链路是否存活
	// 串口设备即便拔掉，Handler 有时也不会立刻返回 error，通过请求测试最可靠
	_, _, _, err := m.Client.ReadCoils(0, 1)
	if err != nil {
		// 先彻底关闭旧句柄
		m.Close()
		// 重新初始化连接
		return m.Connect()
	}
	return nil
}

func (m *ModbusRTUClient) Read(deviceCode, tagCode string, cmdCode int, startAddr int, readLength int, dt string, display string) (any, error) {
	// 自动检查重连
	if err := m.Reconnect(); err != nil {
		return nil, fmt.Errorf("reconnect failed: %v", err)
	}
	switch cmdCode {
	case 0:
		{
			_, _, raw, err := m.Client.ReadCoils(uint16(startAddr-1), 1)
			if err != nil {
				return nil, fmt.Errorf("寄存器(0x-%d)返回数据量错误:%s", startAddr, err.Error())
			}
			if len(raw) != 1 {
				return nil, fmt.Errorf("寄存器(0x-%d)返回数据量错误", startAddr)
			}
			return raw[0], nil
		}
	case 1:
		{
			_, _, raw, err := m.Client.ReadDiscreteInputs(uint16(startAddr-1), 1)
			if err != nil {
				return nil, fmt.Errorf("寄存器(1x-%d)返回数据量错误:%s", startAddr, err.Error())
			}
			if len(raw) != 1 {
				return nil, fmt.Errorf("寄存器(1x-%d)返回数据量错误", startAddr)
			}
			return raw[0], nil
		}
	case 3:
		{
			_, _, raw, err := m.Client.ReadInputRegisters(uint16(startAddr-1), uint16(readLength))
			if err != nil {
				return nil, fmt.Errorf("寄存器(3x-%d-%d)返回数据量错误:%s", startAddr, readLength, err.Error())
			}
			data, err := m.ByteTo(deviceCode, tagCode, raw, dt, display)
			if err != nil {
				return nil, fmt.Errorf("寄存器(3x-%d-%d)返回数据量错误:%s", startAddr, readLength, err.Error())
			}
			return data, nil
		}
	case 4:
		{
			_, _, raw, err := m.Client.ReadHoldingRegisters(uint16(startAddr-1), uint16(readLength))
			if err != nil {
				return nil, fmt.Errorf("寄存器(4x-%d-%d)返回数据量错误:%s", startAddr, readLength, err.Error())
			}
			data, err := m.ByteTo(deviceCode, tagCode, raw, dt, display)
			if err != nil {
				return nil, fmt.Errorf("寄存器(4x-%d-%d)返回数据量错误:%s", startAddr, readLength, err.Error())
			}
			return data, nil
		}
	}
	return nil, fmt.Errorf("错误的指令码:%d", cmdCode)
}
