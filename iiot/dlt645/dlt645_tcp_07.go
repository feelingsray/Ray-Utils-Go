package dlt645

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func NewDLT645TCP07(host string, port int, addr string, timeout int) *DLT645TCP07 {
	dlt := new(DLT645TCP07)
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

type DLT645TCP07 struct {
	Host    string
	Port    int
	Timeout int
	Conn    net.Conn
	Addr    string
	address []byte
}

func (p *DLT645TCP07) CheckConn() error {
	if p.Conn == nil {
		return p.Connect()
	}
	return nil
}

func (p *DLT645TCP07) Connect() error {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", p.Host, p.Port), time.Duration(p.Timeout)*time.Second)
	if err != nil {
		return err
	}
	if conn != nil {
		p.Conn = conn
	}
	return nil
}

func (p *DLT645TCP07) Close() error {
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
func (p *DLT645TCP07) calcCS(data []byte) byte {
	cs := byte(0)
	for _, b := range data {
		cs += b
	}
	return cs
}

// 十进制转Byte
func (p *DLT645TCP07) decimalToBCD(num int) []byte {
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
func (p *DLT645TCP07) strToBCD(data string) ([]byte, error) {
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
func (p *DLT645TCP07) reverseBCD(data []byte) []byte {
	for i, j := 0, len(data)-1; i <= j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

// 解码数据值
func (p *DLT645TCP07) decodeVal(rdData []byte, prec float64) (float64, error) {
	// 数据反转并减去0x33
	for i, j := 0, len(rdData)-1; i <= j; i, j = i+1, j-1 {
		rdData[i], rdData[j] = rdData[j]-0x33, rdData[i]-0x33
	}
	var strVal strings.Builder
	for _, i := range rdData {
		strVal.WriteString(fmt.Sprintf("%02X", i))
	}
	pos := strings.LastIndex(strVal.String(), "F")
	if pos == 16 {
		return 0, nil
	}
	intVal, err := strconv.ParseInt(strVal.String()[pos+1:], 10, 32)
	if err != nil {
		return 0, err
	}
	return float64(intVal) * prec, nil
}

// 监测采集数据
func (p *DLT645TCP07) dlt645TCP07(DI string, prec float64) (float64, error) {
	req := []byte{0x68}                      // 请求头
	req = append(req, p.address...)          // 设备地址
	req = append(req, []byte{0x68, 0x11}...) // 数据域请求头和控制码(11)
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
	if resp[0] != 0x68 || resp[len(resp)-1] != 0x16 || len(resp) < 14 {
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
	rdData := rd[6 : len(rd)-2]
	val, err := p.decodeVal(rdData, prec)
	return val, nil
}

func (p *DLT645TCP07) ReadTag(DI string, prec float64) (float64, error) {
	if p.Conn == nil {
		err := p.Connect()
		if err != nil {
			return 0, err
		}
	}
	defer p.Close()
	return p.dlt645TCP07(DI, prec)
}

func (p *DLT645TCP07) ReadTags(DIP map[string]float64) (map[string]float64, error) {
	if p.Conn == nil {
		err := p.Connect()
		if err != nil {
			return nil, err
		}
	}
	defer p.Close()
	valList := make(map[string]float64)
	for di, prec := range DIP {
		val, err := p.dlt645TCP07(di, prec)
		if err != nil {
			continue
		} else {
			valList[di] = val
		}
		time.Sleep(100 * time.Millisecond)
	}
	return valList, nil
}

// // SetCycleFreeze 设置周期定时冻结
// // minutes: 冻结周期（单位：分钟），传5即设置5分钟周期冻结
// func (p *DLT645TCP07) SetCycleFreeze(minutes int) error {
//	if p.Conn == nil {
//		err := p.Connect()
//		if err != nil {
//			return err
//		}
//	}
//	defer p.Close()
//
//	// 构造4字节冻结命令数据域：周期分钟数(4位BCD) + 0000
//	freezeDataStr := fmt.Sprintf("%04d0000", minutes)
//	bcdR, err := p.strToBCD(freezeDataStr)
//	if err != nil {
//		return err
//	}
//	bcd := p.reverseBCD(bcdR)
//
//	// 数据域加0x33偏移
//	di := make([]byte, 0)
//	for _, d := range bcd {
//		di = append(di, d+0x33)
//	}
//
//	// 组装冻结命令帧（控制码0x16）
//	req := []byte{0x68}
//	req = append(req, p.address...)
//	req = append(req, []byte{0x68, 0x16}...)
//	req = append(req, byte(len(di)))
//	req = append(req, di...)
//	req = append(req, p.calcCS(req))
//	req = append(req, 0x16)
//
//	// 发送命令
//	_, err = p.Conn.Write(req)
//	if err != nil {
//		return err
//	}
//
//	// 接收应答
//	resp := make([]byte, 0)
//	rt := make([]byte, 256)
//	ri, err := p.Conn.Read(rt)
//	if err != nil {
//		return err
//	}
//	for _, r := range rt[:ri] {
//		if r != 0xFE {
//			resp = append(resp, r)
//		}
//	}
//
//	// 帧校验
//	if len(resp) == 0 || resp[0] != 0x68 || resp[len(resp)-1] != 0x16 || len(resp) < 12 {
//		return fmt.Errorf("冻结命令应答帧错误:%X", resp)
//	}
//
//	// 解析应答控制码
//	rdSimple := resp
//	rd := rdSimple[len(p.address)+2:]
//	ctrlCode := rd[0]
//	if ctrlCode == 0xD6 {
//		return fmt.Errorf("冻结命令执行失败，错误码:%X", rd[2])
//	}
//	if ctrlCode != 0x96 {
//		return fmt.Errorf("冻结命令未知应答控制码:%X", ctrlCode)
//	}
//
//	return nil
// }
//
// // 内部读取冻结数据通用方法
// func (p *DLT645TCP07) readFreezeData(DI string, timeStr string, prec float64) (float64, error) {
//	// 拼接数据域：4字节DI + 4字节时间标志(MMDDhhmm)
//	dataStr := DI + timeStr
//	ditR, err := p.strToBCD(dataStr)
//	dit := p.reverseBCD(ditR)
//	if err != nil {
//		return 0, err
//	}
//
//	// 数据域加0x33偏移
//	di := make([]byte, 0)
//	for _, d := range dit {
//		di = append(di, d+0x33)
//	}
//
//	// 组装读数据帧（控制码0x11）
//	req := []byte{0x68}
//	req = append(req, p.address...)
//	req = append(req, []byte{0x68, 0x11}...)
//	req = append(req, byte(len(di)))
//	req = append(req, di...)
//	req = append(req, p.calcCS(req))
//	req = append(req, 0x16)
//
//	// 发送请求
//	_, err = p.Conn.Write(req)
//	if err != nil {
//		return 0, err
//	}
//
//	// 接收应答
//	resp := make([]byte, 0)
//	rt := make([]byte, 4096)
//	ri, err := p.Conn.Read(rt)
//	if err != nil {
//		return 0, err
//	}
//	for _, r := range rt[:ri] {
//		if r != 0xFE {
//			resp = append(resp, r)
//		}
//	}
//
//	// 基本帧校验
//	if len(resp) == 0 {
//		return 0, fmt.Errorf("报文解析错误:%X", resp)
//	}
//	if resp[0] != 0x68 || resp[len(resp)-1] != 0x16 || len(resp) < 14 {
//		return 0, fmt.Errorf("报文解析错误:%X", resp)
//	}
//
//	// 提取有效数据
//	rdSimple := make([]byte, 0)
//	for _, r := range resp {
//		if r == 0xFE {
//			continue
//		} else {
//			rdSimple = append(rdSimple, r)
//		}
//	}
//	rd := rdSimple[len(p.address)+2:]
//	_ = rd[0]                   // 应答功能码
//	rdData := rd[6 : len(rd)-2] // 跳过DI和帧头，取数据值部分
//
//	// 解码数值
//	val, err := p.decodeVal(rdData, prec)
//	return val, nil
// }
//
// // Read5MinFreeze 读取单个时间点的5分钟冻结数据
// // DI: 数据标识（如"00000000"正向有功总电能）
// // t: 目标时间点（精确到分钟）
// // prec: 数值精度（如0.01）
// func (p *DLT645TCP07) Read5MinFreeze(DI string, t time.Time, prec float64) (float64, error) {
//	if p.Conn == nil {
//		err := p.Connect()
//		if err != nil {
//			return 0, err
//		}
//	}
//	defer p.Close()
//
//	// 时间转MMDDhhmm格式BCD字符串
//	timeStr := fmt.Sprintf("%02d%02d%02d%02d", t.Month(), t.Day(), t.Hour(), t.Minute())
//	return p.readFreezeData(DI, timeStr, prec)
// }
//
// // BatchRead5MinFreeze 批量读取指定时间范围内的全部5分钟冻结数据
// // DI: 数据标识
// // startTime: 起始时间
// // endTime: 结束时间
// // prec: 数值精度
// func (p *DLT645TCP07) BatchRead5MinFreeze(DI string, startTime, endTime time.Time, prec float64) (map[time.Time]float64, error) {
//	if p.Conn == nil {
//		err := p.Connect()
//		if err != nil {
//			return nil, err
//		}
//	}
//	defer p.Close()
//
//	result := make(map[time.Time]float64)
//
//	// 按5分钟步长遍历所有时间点
//	for t := startTime; !t.After(endTime); t = t.Add(5 * time.Minute) {
//		timeStr := fmt.Sprintf("%02d%02d%02d%02d", t.Month(), t.Day(), t.Hour(), t.Minute())
//		val, err := p.readFreezeData(DI, timeStr, prec)
//		if err != nil {
//			continue // 单个时间点失败不中断整体流程
//		}
//		result[t] = val
//		time.Sleep(100 * time.Millisecond) // 与原批量读取保持相同间隔
//	}
//
//	return result, nil
// }
