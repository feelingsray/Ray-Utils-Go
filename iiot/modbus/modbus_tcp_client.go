package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/ivanlebron/modbus"
	"github.com/spf13/cast"
)

func NewModbusTCPClient(host string, slaverId int, timeout int) *ModbusTCPClient {
	mc := new(ModbusTCPClient)
	mc.Host = host
	mc.SlaveId = slaverId
	mc.Timeout = timeout
	return mc
}

type ModbusTCPClient struct {
	Host       string `json:"host"`
	SlaveId    int    `json:"slave_id"`
	Timeout    int    `json:"timeout"`
	TCPHandler *modbus.TCPClientHandler
	Client     modbus.Client
}

func (m *ModbusTCPClient) LinerScalingSimple(value float64, scaling float64) (float64, error) {
	data := value * scaling
	return data, nil
}

func (m *ModbusTCPClient) LinerScaling(value any, offset float64, rawHeight float64, rawLow float64, scaleHeight float64, scaleLow float64) (float64, error) {
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

func (m *ModbusTCPClient) ByteTo(data []byte, dataType, displayType string) (any, error) {
	switch dataType {
	case "float":
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
	case "double":
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
	case "long", "int32":
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
	case "int64":
		if len(data) != 8 {
			return nil, errors.New("数据长度不正确,应为8个8位")
		}
		if displayType == "ABCDEFGH" {
			bits := binary.BigEndian.Uint64(data)
			return int64(bits), nil
		}
		if displayType == "GHEFCDAB" {
			bits := binary.BigEndian.Uint64([]byte{data[6], data[7], data[4], data[5], data[2], data[3], data[0], data[1]})
			return int64(bits), nil
		}
		if displayType == "BADCFEHG" {
			bits := binary.BigEndian.Uint64([]byte{data[1], data[0], data[3], data[2], data[5], data[4], data[7], data[6]})
			return int64(bits), nil
		}
		if displayType == "HGFEDCBA" {
			bits := binary.BigEndian.Uint64([]byte{data[7], data[6], data[5], data[4], data[3], data[2], data[1], data[0]})
			return int64(bits), nil
		}
	case "int", "int16":
		if len(data) != 2 {
			return nil, errors.New("数据长度不正确,应为2个8位")
		}
		if displayType == "AB" {
			bits := binary.BigEndian.Uint16(data)
			return int16(bits), nil
		}
		bits := binary.LittleEndian.Uint16(data)
		return int16(bits), nil
	// ========== 新增 BCD 支持 ==========
	case "bcd", "BCD":
		switch len(data) {
		case 2:
			var val uint16
			if displayType == "AB" { // 大端
				val = binary.BigEndian.Uint16(data)
			} else { // 默认小端（如 "BA"）
				val = binary.LittleEndian.Uint16(data)
			}
			// 可选：校验每个半字节是否为 0~9，此处略
			return uint64(val), nil
		case 4:
			var val uint32
			if displayType == "ABCD" { // 大端
				val = binary.BigEndian.Uint32(data)
			} else { // 默认小端（如 "DCBA"）
				val = binary.LittleEndian.Uint32(data)
			}
			return uint64(val), nil
		case 8:
			var val uint64
			if displayType == "ABCDEFGH" { // 大端
				val = binary.BigEndian.Uint64(data)
			} else { // 默认小端（如 "HGFEDCBA"）
				val = binary.LittleEndian.Uint64(data)
			}
			return val, nil
		default:
			return nil, errors.New("BCD 数据长度只支持 2、4 或 8 字节")
		}
	// =================================
	default:
		return nil, errors.New("data type not support")
	}
	return nil, nil
}

func (m *ModbusTCPClient) ToByte(data any, dataType string, displayType string) ([]byte, error) {
	switch dataType {
	case "float":
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
	case "long", "int32":
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
	case "double":
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
	case "int64":
		if displayType == "ABCDEFGH" {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint64(bytes, uint64(data.(int32)))
			return bytes, nil
		}
		if displayType == "GHEFCDAB" {
			bits := uint64(data.(int32))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return []byte{bytes[6], bytes[7], bytes[4], bytes[5], bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		}
		if displayType == "BADCFEHG" {
			bits := uint64(data.(int32))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2], bytes[5], bytes[4], bytes[7], bytes[6]}, nil
		}
		if displayType == "HGFEDCBA" {
			bits := uint64(data.(int32))
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, bits)
			return []byte{bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		}
	// ========== 新增 BCD 编码支持 ==========
	case "bcd", "BCD":
		// 将 data 转为 uint64（支持 uint64, uint32, uint16, int 等）
		var val uint64
		switch v := data.(type) {
		case uint64:
			val = v
		case uint32:
			val = uint64(v)
		case uint16:
			val = uint64(v)
		case int:
			if v < 0 {
				return nil, errors.New("BCD 不支持负数")
			}
			val = uint64(v)
		default:
			return nil, errors.New("BCD 数据必须为无符号整数类型")
		}
		// 根据 displayType 推断字节长度
		var length int
		switch displayType {
		case "ABCD", "DCBA", "CDAB", "BADC":
			length = 4
		case "ABCDEFGH", "HGFEDCBA", "GHEFCDAB", "BADCFEHG":
			length = 8
		default: // "AB", "BA" 等
			length = 2
		}
		// 检查数值是否超出 BCD 可表示范围
		maxVal := uint64(0)
		switch length {
		case 2:
			maxVal = 9999
		case 4:
			maxVal = 99999999
		case 8:
			maxVal = 9999999999999999
		default:
			return nil, errors.New("不支持的 BCD 长度")
		}
		if val > maxVal {
			return nil, fmt.Errorf("BCD 数值 %d 超出 %d 字节表示范围", val, length)
		}
		// 将 val 格式化为十进制字符串，不足位补零
		str := fmt.Sprintf("%0*d", length*2, val)
		bytes := make([]byte, length)
		for i := 0; i < length; i++ {
			high := str[i*2] - '0'
			low := str[i*2+1] - '0'
			bytes[i] = (high << 4) | low
		}
		// 根据 displayType 处理字节序（大端/小端）
		// 大端标识：2字节 "AB"，4字节 "ABCD"，8字节 "ABCDEFGH"
		// 否则视为小端
		if length == 2 && displayType != "AB" {
			bytes[0], bytes[1] = bytes[1], bytes[0]
		} else if length == 4 && displayType != "ABCD" {
			for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
				bytes[i], bytes[j] = bytes[j], bytes[i]
			}
		} else if length == 8 && displayType != "ABCDEFGH" {
			for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
				bytes[i], bytes[j] = bytes[j], bytes[i]
			}
		}
		return bytes, nil
	// =====================================
	default:
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

func (m *ModbusTCPClient) CheckConn() error {
	if m.Client == nil {
		return m.Connect()
	}
	return nil
}

func (m *ModbusTCPClient) Connect() error {
	handler := modbus.NewTCPClientHandler(m.Host)
	handler.SlaveId = byte(m.SlaveId)
	handler.Timeout = time.Duration(m.Timeout) * time.Second
	handler.IdleTimeout = time.Duration(m.Timeout) * time.Second
	err := handler.Connect()
	if err != nil {
		m.Client = nil
		return err
	}
	m.TCPHandler = handler
	m.Client = modbus.NewClient(handler)
	return nil
}

func (m *ModbusTCPClient) Close() error {
	if m.Client != nil {
		m.Client = nil
	}
	if m.TCPHandler != nil {
		return m.TCPHandler.Close()
	}
	return nil
}

func (m *ModbusTCPClient) Read(cmdCode, startAddr, readLength int, dt, display string) (ReadRes, error) {
	var res ReadRes
	if m.Client == nil {
		return res, errors.New("client nil")
	}
	switch cmdCode {
	case 0:
		{
			w, r, raw, err := m.Client.ReadCoils(uint16(startAddr-1), 1)
			if err != nil {
				return res, fmt.Errorf("寄存器(0x-%d)读取错误:%s", startAddr, err.Error())
			}
			if len(raw) != 1 {
				return res, fmt.Errorf("寄存器(0x-%d)返回数据量错误", startAddr)
			}
			res.Write, res.Read, res.Val = w, r, cast.ToFloat64(raw[0])
			return res, nil
		}
	case 1:
		{
			w, r, raw, err := m.Client.ReadDiscreteInputs(uint16(startAddr-1), 1)
			if err != nil {
				return res, fmt.Errorf("寄存器(1x-%d)读取错误:%s", startAddr, err.Error())
			}
			if len(raw) != 1 {
				return res, fmt.Errorf("寄存器(1x-%d)返回数据量错误", startAddr)
			}
			res.Write, res.Read, res.Val = w, r, cast.ToFloat64(raw[0])
			return res, nil
		}
	case 3:
		{
			w, r, raw, err := m.Client.ReadInputRegisters(uint16(startAddr-1), uint16(readLength))
			if err != nil {
				return res, fmt.Errorf("寄存器(3x-%d-%d)读取错误:%s", startAddr, readLength, err.Error())
			}
			data, err := m.ByteTo(raw, dt, display)
			if err != nil {
				return res, fmt.Errorf("寄存器(3x-%d-%d)返回数据量错误:%s", startAddr, readLength, err.Error())
			}
			res.Write, res.Read, res.Val = w, r, cast.ToFloat64(data)
			return res, nil
		}
	case 4:
		{
			w, r, raw, err := m.Client.ReadHoldingRegisters(uint16(startAddr-1), uint16(readLength))
			if err != nil {
				return res, fmt.Errorf("寄存器(4x-%d-%d)读取错误:%s", startAddr, readLength, err.Error())
			}
			data, err := m.ByteTo(raw, dt, display)
			if err != nil {
				return res, fmt.Errorf("寄存器(4x-%d-%d)返回数据量错误:%s", startAddr, readLength, err.Error())
			}
			res.Write, res.Read, res.Val = w, r, cast.ToFloat64(data)
			return res, nil
		}
	}
	return res, fmt.Errorf("错误的指令码:%d", cmdCode)
}
