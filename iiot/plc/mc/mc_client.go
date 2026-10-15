package mc

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"github.com/future-architect/go-mcprotocol/mcp"
)

// 三菱 PLC

const defaultMcPort = 5007

var bitDevices = map[string]struct{}{
	"X": {}, "Y": {}, "M": {}, "L": {}, "F": {}, "V": {}, "B": {},
}

func NewMCClient(host string, port int, timeout int, byteOrder string, networkNum, pcNum, unitIoNum, unitStationNum string) *MCClient {
	return &MCClient{
		Host:           host,
		Port:           port,
		Timeout:        timeout,
		ByteOrder:      byteOrder,
		NetworkNum:     networkNum,
		PcNum:          pcNum,
		UnitIoNum:      unitIoNum,
		UnitStationNum: unitStationNum,
	}
}

type MCClient struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Timeout        int    `json:"timeout"`
	ByteOrder      string `json:"byte_order"`
	NetworkNum     string `json:"network_num"`
	PcNum          string `json:"pc_num"`
	UnitIoNum      string `json:"unit_io_num"`
	UnitStationNum string `json:"unit_station_num"`
	Client         mcp.Client
	wordSwap       bool
}

func (m *MCClient) LinerScalingSimple(value float64, scaling float64) (float64, error) {
	return value * scaling, nil
}

func (m *MCClient) Connect() error {
	port := m.Port
	if port == 0 {
		port = defaultMcPort
	}
	timeout := time.Duration(m.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", m.Host, port), timeout)
	if err != nil {
		return err
	}
	_ = conn.Close()
	var client mcp.Client
	if m.NetworkNum == "" && m.PcNum == "" && m.UnitIoNum == "" && m.UnitStationNum == "" {
		client, err = mcp.New3EClient(m.Host, port, mcp.NewLocalStation())
	} else {
		networkNum := m.NetworkNum
		if networkNum == "" {
			networkNum = "00"
		}
		pcNum := m.PcNum
		if pcNum == "" {
			pcNum = "FF"
		}
		unitIoNum := m.UnitIoNum
		if unitIoNum == "" {
			unitIoNum = "FF03"
		}
		unitStationNum := m.UnitStationNum
		if unitStationNum == "" {
			unitStationNum = "00"
		}
		client, err = mcp.New3EClient(m.Host, port, mcp.NewStation(networkNum, pcNum, unitIoNum, unitStationNum))
	}
	if err != nil {
		return err
	}
	m.wordSwap = strings.ToLower(m.ByteOrder) == "big"
	m.Client = client
	return nil
}

func (m *MCClient) uint32FromWords(raw []byte) uint32 {
	w0 := binary.LittleEndian.Uint16(raw[0:2])
	w1 := binary.LittleEndian.Uint16(raw[2:4])
	if m.wordSwap {
		return uint32(w0)<<16 | uint32(w1)
	}
	return uint32(w1)<<16 | uint32(w0)
}

func (m *MCClient) uint64FromWords(raw []byte) uint64 {
	w0 := uint64(binary.LittleEndian.Uint16(raw[0:2]))
	w1 := uint64(binary.LittleEndian.Uint16(raw[2:4]))
	w2 := uint64(binary.LittleEndian.Uint16(raw[4:6]))
	w3 := uint64(binary.LittleEndian.Uint16(raw[6:8]))
	if m.wordSwap {
		return w0<<48 | w1<<32 | w2<<16 | w3
	}
	return w3<<48 | w2<<32 | w1<<16 | w0
}

func (m *MCClient) Close() error {
	m.Client = nil
	return nil
}

func (m *MCClient) Read(deviceCode, tagCode, device string, startAddr, bitOffset, readLength int, dt string) (any, error) {
	device = strings.ToUpper(strings.TrimSpace(device))
	if device == "" {
		return nil, fmt.Errorf("device %s, tag %s, 设备类型未设置", deviceCode, tagCode)
	}
	switch strings.ToLower(dt) {
	case "bool", "bit":
		return m.readBit(device, startAddr, bitOffset)
	case "byte", "uint8":
		raw, err := m.readWords(device, startAddr, 1)
		if err != nil {
			return nil, err
		}
		return raw[0], nil
	case "bytes":
		if readLength <= 0 {
			return nil, fmt.Errorf("设备(%s-%d)读取长度错误", device, startAddr)
		}
		points := (readLength + 1) / 2
		raw, err := m.readWords(device, startAddr, points)
		if err != nil {
			return nil, err
		}
		if len(raw) < readLength {
			return nil, fmt.Errorf("设备(%s-%d)返回数据量错误", device, startAddr)
		}
		return raw[:readLength], nil
	case "string":
		if readLength <= 0 {
			return nil, fmt.Errorf("设备(%s-%d)读取长度错误", device, startAddr)
		}
		points := (readLength + 1) / 2
		raw, err := m.readWords(device, startAddr, points)
		if err != nil {
			return nil, err
		}
		if len(raw) < readLength {
			return nil, fmt.Errorf("设备(%s-%d)返回数据量错误", device, startAddr)
		}
		return strings.TrimRight(string(raw[:readLength]), "\x00"), nil
	case "word", "uint16":
		raw, err := m.readWords(device, startAddr, 1)
		if err != nil {
			return nil, err
		}
		return binary.LittleEndian.Uint16(raw), nil
	case "int", "int16":
		raw, err := m.readWords(device, startAddr, 1)
		if err != nil {
			return nil, err
		}
		return int16(binary.LittleEndian.Uint16(raw)), nil
	case "dword", "uint32", "udint":
		raw, err := m.readWords(device, startAddr, 2)
		if err != nil {
			return nil, err
		}
		return m.uint32FromWords(raw), nil
	case "dint", "int32":
		raw, err := m.readWords(device, startAddr, 2)
		if err != nil {
			return nil, err
		}
		return int32(m.uint32FromWords(raw)), nil
	case "real", "float", "float32":
		raw, err := m.readWords(device, startAddr, 2)
		if err != nil {
			return nil, err
		}
		return math.Float32frombits(m.uint32FromWords(raw)), nil
	case "lword", "uint64", "ulint":
		raw, err := m.readWords(device, startAddr, 4)
		if err != nil {
			return nil, err
		}
		return m.uint64FromWords(raw), nil
	case "lint", "int64":
		raw, err := m.readWords(device, startAddr, 4)
		if err != nil {
			return nil, err
		}
		return int64(m.uint64FromWords(raw)), nil
	case "lreal", "double", "float64":
		raw, err := m.readWords(device, startAddr, 4)
		if err != nil {
			return nil, err
		}
		return math.Float64frombits(m.uint64FromWords(raw)), nil
	default:
		return nil, fmt.Errorf("device %s, tag %s, 未知数据类型:%s", deviceCode, tagCode, dt)
	}
}

func (m *MCClient) readBit(device string, startAddr, bitOffset int) (bool, error) {
	if _, ok := bitDevices[device]; ok {
		resp, err := m.Client.BitRead(device, int64(startAddr), 1)
		if err != nil {
			return false, fmt.Errorf("设备(%s-%d)读取失败:%s", device, startAddr, err.Error())
		}
		payload, err := m.parsePayload(resp, device, startAddr)
		if err != nil {
			return false, err
		}
		if len(payload) < 1 {
			return false, fmt.Errorf("设备(%s-%d)返回数据量错误", device, startAddr)
		}
		return (payload[0] & 0x10) != 0, nil
	}
	raw, err := m.readWords(device, startAddr, 1)
	if err != nil {
		return false, err
	}
	word := binary.LittleEndian.Uint16(raw)
	if bitOffset < 0 || bitOffset > 15 {
		return false, fmt.Errorf("设备(%s-%d)位偏移错误:%d", device, startAddr, bitOffset)
	}
	return (word>>uint(bitOffset))&0x01 == 1, nil
}

func (m *MCClient) readWords(device string, startAddr, numPoints int) ([]byte, error) {
	if numPoints <= 0 {
		return nil, fmt.Errorf("设备(%s-%d)读取长度错误", device, startAddr)
	}
	resp, err := m.Client.Read(device, int64(startAddr), int64(numPoints))
	if err != nil {
		return nil, fmt.Errorf("设备(%s-%d-%d)读取失败:%s", device, startAddr, numPoints, err.Error())
	}
	payload, err := m.parsePayload(resp, device, startAddr)
	if err != nil {
		return nil, err
	}
	if len(payload) < 2*numPoints {
		return nil, fmt.Errorf("设备(%s-%d)返回数据量错误", device, startAddr)
	}
	return payload[:2*numPoints], nil
}

func (m *MCClient) parsePayload(resp []byte, device string, startAddr int) ([]byte, error) {
	parsed, err := mcp.NewParser().Do(resp)
	if err != nil {
		return nil, fmt.Errorf("设备(%s-%d)响应解析失败:%s", device, startAddr, err.Error())
	}
	if parsed.EndCode != "0000" {
		return nil, fmt.Errorf("设备(%s-%d)响应错误码:%s", device, startAddr, parsed.EndCode)
	}
	return parsed.Payload, nil
}
