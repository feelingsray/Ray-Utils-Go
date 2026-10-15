package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ivanlebron/modbus"
	"github.com/spf13/cast"
)

func NewModbusTCPClientPro(host string, slaverId int, timeout int) *ModbusTCPClientPro {
	mc := new(ModbusTCPClientPro)
	mc.Host = host
	mc.SlaveId = slaverId
	mc.Timeout = timeout
	return mc
}

type ModbusTCPClientPro struct {
	Host       string `json:"host"`
	SlaveId    int    `json:"slave_id"`
	Timeout    int    `json:"timeout"`
	TCPHandler *modbus.TCPClientHandler
	Client     modbus.Client
	mu         sync.Mutex
}

func (m *ModbusTCPClientPro) LinerScalingSimple(value float64, scaling float64) (float64, error) {
	data := value * scaling
	return data, nil
}

func (m *ModbusTCPClientPro) LinerScaling(value any, offset float64, rawHeight float64, rawLow float64, scaleHeight float64, scaleLow float64) (float64, error) {
	if value == nil {
		return 0, errors.New("value 不能为空")
	}
	v, ok := value.(float64)
	if !ok {
		return 0, errors.New("value 必须是 float64 类型")
	}
	if rawHeight == rawLow {
		return 0, errors.New("rawHeight 与 rawLow 不能相等")
	}
	data := (v-offset)/(rawHeight-rawLow)*(scaleHeight-scaleLow) + scaleLow
	return data, nil
}

func (m *ModbusTCPClientPro) ByteTo(data []byte, dataType, displayType string) (any, error) {
	switch dataType {
	case "float":
		if len(data) != 4 {
			return nil, errors.New("数据长度不正确,应为4个8位")
		}
		switch displayType {
		case "ABCD":
			bits := binary.BigEndian.Uint32(data)
			return math.Float32frombits(bits), nil
		case "CDAB":
			bits := binary.BigEndian.Uint32([]byte{data[2], data[3], data[0], data[1]})
			return math.Float32frombits(bits), nil
		case "DCBA":
			bits := binary.BigEndian.Uint32([]byte{data[3], data[2], data[1], data[0]})
			return math.Float32frombits(bits), nil
		case "BADC":
			bits := binary.BigEndian.Uint32([]byte{data[1], data[0], data[3], data[2]})
			return math.Float32frombits(bits), nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "double":
		if len(data) != 8 {
			return nil, errors.New("数据长度不正确,应为8个8位")
		}
		switch displayType {
		case "ABCDEFGH":
			bits := binary.BigEndian.Uint64(data)
			return math.Float64frombits(bits), nil
		case "GHEFCDAB":
			bits := binary.BigEndian.Uint64([]byte{data[6], data[7], data[4], data[5], data[2], data[3], data[0], data[1]})
			return math.Float64frombits(bits), nil
		case "BADCFEHG":
			bits := binary.BigEndian.Uint64([]byte{data[1], data[0], data[3], data[2], data[5], data[4], data[7], data[6]})
			return math.Float64frombits(bits), nil
		case "HGFEDCBA":
			bits := binary.BigEndian.Uint64([]byte{data[7], data[6], data[5], data[4], data[3], data[2], data[1], data[0]})
			return math.Float64frombits(bits), nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "long", "int32":
		if len(data) != 4 {
			return nil, errors.New("数据长度不正确,应为4个8位")
		}
		switch displayType {
		case "ABCD":
			bits := binary.BigEndian.Uint32(data)
			return int32(bits), nil
		case "CDAB":
			bits := binary.BigEndian.Uint32([]byte{data[2], data[3], data[0], data[1]})
			return int32(bits), nil
		case "DCBA":
			bits := binary.BigEndian.Uint32([]byte{data[3], data[2], data[1], data[0]})
			return int32(bits), nil
		case "BADC":
			bits := binary.BigEndian.Uint32([]byte{data[1], data[0], data[3], data[2]})
			return int32(bits), nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "int64":
		if len(data) != 8 {
			return nil, errors.New("数据长度不正确,应为8个8位")
		}
		switch displayType {
		case "ABCDEFGH":
			bits := binary.BigEndian.Uint64(data)
			return int64(bits), nil
		case "GHEFCDAB":
			bits := binary.BigEndian.Uint64([]byte{data[6], data[7], data[4], data[5], data[2], data[3], data[0], data[1]})
			return int64(bits), nil
		case "BADCFEHG":
			bits := binary.BigEndian.Uint64([]byte{data[1], data[0], data[3], data[2], data[5], data[4], data[7], data[6]})
			return int64(bits), nil
		case "HGFEDCBA":
			bits := binary.BigEndian.Uint64([]byte{data[7], data[6], data[5], data[4], data[3], data[2], data[1], data[0]})
			return int64(bits), nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "int", "int16":
		if len(data) != 2 {
			return nil, errors.New("数据长度不正确,应为2个8位")
		}
		if displayType == "AB" {
			bits := binary.BigEndian.Uint16(data)
			return int16(bits), nil
		} else if displayType == "BA" {
			bits := binary.LittleEndian.Uint16(data)
			return int16(bits), nil
		} else {
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "bcd", "BCD":
		var byteLen int
		switch len(data) {
		case 2:
			byteLen = 2
		case 4:
			byteLen = 4
		case 8:
			byteLen = 8
		default:
			return nil, errors.New("BCD 数据长度只支持 2、4 或 8 字节")
		}
		bytes := make([]byte, byteLen)
		copy(bytes, data)
		isBigEndian := false
		switch byteLen {
		case 2:
			isBigEndian = (displayType == "AB")
		case 4:
			isBigEndian = (displayType == "ABCD")
		case 8:
			isBigEndian = (displayType == "ABCDEFGH")
		}
		if !isBigEndian {
			for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
				bytes[i], bytes[j] = bytes[j], bytes[i]
			}
		}
		var val uint64
		for _, b := range bytes {
			high := b >> 4
			low := b & 0x0F
			if high > 9 || low > 9 {
				return nil, fmt.Errorf("无效的 BCD 字节: 0x%02X", b)
			}
			val = val*100 + uint64(high)*10 + uint64(low)
		}
		return val, nil
	default:
		return nil, errors.New("data type not support")
	}
}

func (m *ModbusTCPClientPro) ToByte(data any, dataType string, displayType string) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data 不能为空")
	}
	switch dataType {
	case "float":
		v, ok := data.(float32)
		if !ok {
			return nil, errors.New("数据类型不匹配，需要 float32")
		}
		bits := math.Float32bits(v)
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits)
		switch displayType {
		case "ABCD":
			return bytes, nil
		case "CDAB":
			return []byte{bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		case "DCBA":
			return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		case "BADC":
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2]}, nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "long", "int32":
		v, ok := data.(int32)
		if !ok {
			return nil, errors.New("数据类型不匹配，需要 int32")
		}
		bits := uint32(v)
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits)
		switch displayType {
		case "ABCD":
			return bytes, nil
		case "CDAB":
			return []byte{bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		case "DCBA":
			return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		case "BADC":
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2]}, nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "double":
		v, ok := data.(float64)
		if !ok {
			return nil, errors.New("数据类型不匹配，需要 float64")
		}
		bits := math.Float64bits(v)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits)
		switch displayType {
		case "ABCDEFGH":
			return bytes, nil
		case "GHEFCDAB":
			return []byte{bytes[6], bytes[7], bytes[4], bytes[5], bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		case "BADCFEHG":
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2], bytes[5], bytes[4], bytes[7], bytes[6]}, nil
		case "HGFEDCBA":
			return []byte{bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "int64":
		v, ok := data.(int64)
		if !ok {
			return nil, errors.New("数据类型不匹配，需要 int64")
		}
		bits := uint64(v)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits)
		switch displayType {
		case "ABCDEFGH":
			return bytes, nil
		case "GHEFCDAB":
			return []byte{bytes[6], bytes[7], bytes[4], bytes[5], bytes[2], bytes[3], bytes[0], bytes[1]}, nil
		case "BADCFEHG":
			return []byte{bytes[1], bytes[0], bytes[3], bytes[2], bytes[5], bytes[4], bytes[7], bytes[6]}, nil
		case "HGFEDCBA":
			return []byte{bytes[7], bytes[6], bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]}, nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "int", "int16":
		v, ok := data.(int16)
		if !ok {
			return nil, errors.New("数据类型不匹配，需要 int16")
		}
		bits := uint16(v)
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, bits)
		switch displayType {
		case "AB":
			return bytes, nil
		case "BA":
			return []byte{bytes[1], bytes[0]}, nil
		default:
			return nil, fmt.Errorf("不支持的字节序: %s", displayType)
		}
	case "bcd", "BCD":
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
		case int64:
			if v < 0 {
				return nil, errors.New("BCD 不支持负数")
			}
			val = uint64(v)
		default:
			return nil, errors.New("BCD 数据必须为无符号整数类型")
		}
		var length int
		switch displayType {
		case "ABCD", "DCBA", "CDAB", "BADC":
			length = 4
		case "ABCDEFGH", "HGFEDCBA", "GHEFCDAB", "BADCFEHG":
			length = 8
		default:
			length = 2
		}
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
		str := fmt.Sprintf("%0*d", length*2, val)
		bytes := make([]byte, length)
		for i := 0; i < length; i++ {
			high := str[i*2] - '0'
			low := str[i*2+1] - '0'
			bytes[i] = (high << 4) | low
		}
		isBigEndian := false
		switch length {
		case 2:
			isBigEndian = (displayType == "AB")
		case 4:
			isBigEndian = (displayType == "ABCD")
		case 8:
			isBigEndian = (displayType == "ABCDEFGH")
		}
		if !isBigEndian {
			for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
				bytes[i], bytes[j] = bytes[j], bytes[i]
			}
		}
		return bytes, nil
	default:
		return nil, errors.New("data type not support")
	}
}

func (m *ModbusTCPClientPro) CheckConn() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Client == nil {
		return m.connect()
	}
	return nil
}

func (m *ModbusTCPClientPro) connect() error {
	if m.SlaveId < 1 || m.SlaveId > 247 {
		return fmt.Errorf("SlaveId 必须在 1-247 之间，当前: %d", m.SlaveId)
	}
	if m.Timeout <= 0 {
		m.Timeout = 5
	}
	handler := modbus.NewTCPClientHandler(m.Host)
	handler.SlaveId = byte(m.SlaveId)
	handler.Timeout = time.Duration(m.Timeout) * time.Second
	handler.IdleTimeout = time.Duration(m.Timeout) * time.Second
	err := handler.Connect()
	if err != nil {
		m.Client = nil
		m.TCPHandler = nil
		return err
	}
	m.TCPHandler = handler
	m.Client = modbus.NewClient(handler)
	return nil
}

func (m *ModbusTCPClientPro) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.connect()
}

func (m *ModbusTCPClientPro) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Client != nil {
		m.Client = nil
	}
	if m.TCPHandler != nil {
		err := m.TCPHandler.Close()
		m.TCPHandler = nil
		return err
	}
	return nil
}

func (m *ModbusTCPClientPro) Read(cmdCode, startAddr, readLength int, dt, display string) (ReadRes, error) {
	var res ReadRes
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Client == nil {
		if err := m.connect(); err != nil {
			return res, fmt.Errorf("连接失败: %w", err)
		}
	}

	if startAddr < 1 || startAddr > 65536 {
		return res, fmt.Errorf("非法起始地址: %d", startAddr)
	}
	if readLength < 1 {
		return res, fmt.Errorf("非法读取长度: %d", readLength)
	}

	switch cmdCode {
	case 0:
		w, r, raw, err := m.Client.ReadCoils(uint16(startAddr-1), 1)
		if err != nil {
			return res, fmt.Errorf("寄存器(0x-%d)读取错误:%s", startAddr, err.Error())
		}
		if len(raw) != 1 {
			return res, fmt.Errorf("寄存器(0x-%d)返回数据量错误", startAddr)
		}
		res.Write, res.Read, res.Val = w, r, cast.ToFloat64(raw[0])
		return res, nil
	case 1:
		w, r, raw, err := m.Client.ReadDiscreteInputs(uint16(startAddr-1), 1)
		if err != nil {
			return res, fmt.Errorf("寄存器(1x-%d)读取错误:%s", startAddr, err.Error())
		}
		if len(raw) != 1 {
			return res, fmt.Errorf("寄存器(1x-%d)返回数据量错误", startAddr)
		}
		res.Write, res.Read, res.Val = w, r, cast.ToFloat64(raw[0])
		return res, nil
	case 3:
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
	case 4:
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
	return res, fmt.Errorf("错误的指令码:%d", cmdCode)
}
