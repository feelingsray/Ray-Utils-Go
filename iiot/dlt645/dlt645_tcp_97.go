package dlt645

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func NewDLT645TCP97(host string, port int, addr string, timeout int) *DLT645TCP97 {
	dlt := new(DLT645TCP97)
	dlt.Host = host
	dlt.Port = port
	dlt.Timeout = timeout
	dlt.Addr = addr
	address, _ := dlt.strToBCD(addr)
	dlt.address = dlt.reverseBCD(address)
	if dlt.Conn == nil {
		_ = dlt.Connect()
	}
	return dlt
}

type DLT645TCP97 struct {
	Host    string
	Port    int
	Timeout int
	Conn    net.Conn
	Addr    string
	address []byte
}

func (p *DLT645TCP97) CheckConn() error {
	if p.Conn == nil {
		return p.Connect()
	}
	return nil
}

func (p *DLT645TCP97) Connect() error {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", p.Host, p.Port), time.Duration(p.Timeout)*time.Second)
	if err != nil {
		return err
	}
	if conn != nil {
		p.Conn = conn
	}
	return nil
}

func (p *DLT645TCP97) Close() error {
	if p.Conn != nil {
		err := p.Conn.Close()
		if err != nil {
			p.Conn = nil
			return err
		}
		p.Conn = nil
	}
	return nil
}

// 求和校验
func (p *DLT645TCP97) calcCS(data []byte) byte {
	cs := byte(0)
	for _, b := range data {
		cs += b
	}
	return cs
}

// 十进制转Byte
func (p *DLT645TCP97) decimalToBCD(num int) []byte {
	str := strconv.Itoa(num)
	bcd := make([]byte, (len(str)+1)/2)
	for i, j := 0, 0; i < len(str); i, j = i+2, j+1 {
		if i+1 == len(str) {
			bcd[j] = str[i] - '0'
		} else {
			bcd[j] = (str[i]-'0')<<4 | (str[i+1] - '0')
		}
	}
	return bcd
}

// 字符串转Byte
func (p *DLT645TCP97) strToBCD(data string) ([]byte, error) {
	if len(data)%2 != 0 {
		data = "0" + data // 如果字符串长度为奇数，则在前面添加一个0字符
	}
	bcd := make([]byte, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		high, err := strconv.ParseUint(string(data[i]), 16, 4)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseUint(string(data[i+1]), 16, 4)
		if err != nil {
			return nil, err
		}
		bcd[i/2] = byte(high<<4) | byte(low)
	}
	return bcd, nil
}

// 反转
func (p *DLT645TCP97) reverseBCD(data []byte) []byte {
	for i, j := 0, len(data)-1; i <= j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

// 解码数据值
func (p *DLT645TCP97) decodeVal(rdData []byte, prec float64) (float64, error) {
	// 数据反转并减去0x33
	for i, j := 0, len(rdData)-1; i <= j; i, j = i+1, j-1 {
		rdData[i], rdData[j] = rdData[j]-0x33, rdData[i]-0x33
	}
	strVal := ""
	for _, i := range rdData {
		strVal += fmt.Sprintf("%02X", i)
	}
	pos := strings.LastIndex(strVal, "F")
	if pos == 8 {
		return 0, nil
	}
	intVal, err := strconv.ParseInt(strVal[pos+1:], 10, 16)
	if err != nil {
		return 0, err
	}
	return float64(intVal) * prec, nil
}

// 监测采集数据
func (p *DLT645TCP97) dlt645TCP97(DI string, prec float64) (float64, error) {
	req := []byte{0x68}                      // 请求头
	req = append(req, p.address...)          // 设备地址
	req = append(req, []byte{0x68, 0x01}...) // 数据域请求头和控制码(01)
	// 转换数据项编码,并加0x33
	ditR, err := p.strToBCD(DI)
	dit := p.reverseBCD(ditR)
	if err != nil {
		return 0, err
	}
	di := make([]byte, 0)
	for _, d := range dit {
		di = append(di, d+0x33)
	}
	req = append(req, byte(len(di))) // 数据项编码长度
	req = append(req, di...)         // 数据项编码
	req = append(req, p.calcCS(req)) // 校验
	req = append(req, 0x16)          // 结束符
	if p.Conn == nil {
		return 0, errors.New("conn is nil")
	}
	_, err = p.Conn.Write(req)
	if err != nil {
		return 0, err
	}
	resp := make([]byte, 0)
	rt := make([]byte, 4096)
	ri, err := p.Conn.Read(rt)
	if err != nil {
		return 0, err
	}
	for _, r := range rt[:ri] {
		if r != 0xFE {
			resp = append(resp, r)
		}
	}
	if len(resp) == 0 {
		return 0, fmt.Errorf("报文解析错误:%X", resp)
	}
	if (resp[0] != 0x68) || resp[len(resp)-1] != 0x16 || len(resp) < 14 {
		return 0, fmt.Errorf("报文解析错误:%X", resp)
	}
	// 过滤无用信息获取有用信息
	rdSimple := make([]byte, 0)
	// 过滤 0xFE
	for _, r := range resp {
		if r == 0xFE {
			continue
		} else {
			rdSimple = append(rdSimple, r)
		}
	}
	// 过滤 前面的信息
	rd := rdSimple[len(p.address)+2:]
	_ = rd[0] // 应答功能码
	rdData := rd[4 : len(rd)-2]
	val, err := p.decodeVal(rdData, prec)
	return val, nil
}

// ReadTag 一次访问一个值
func (p *DLT645TCP97) ReadTag(DI string, prec float64) (float64, error) {
	if p.Conn == nil {
		err := p.Connect()
		if err != nil {
			return 0, err
		}
	}
	defer p.Close()
	return p.dlt645TCP97(DI, prec)
}

// ReadTags 一次访问多个值
// 入:key-DI,v:prec,出:key-DI,v:值
func (p *DLT645TCP97) ReadTags(DIP map[string]float64) (map[string]float64, error) {
	if p.Conn == nil {
		err := p.Connect()
		if err != nil {
			return nil, err
		}
	}
	defer p.Close()
	valList := make(map[string]float64)
	for di, prec := range DIP {
		val, err := p.dlt645TCP97(di, prec)
		if err != nil {
			continue
		} else {
			valList[di] = val
		}
		time.Sleep(100 * time.Millisecond)
	}
	return valList, nil
}
