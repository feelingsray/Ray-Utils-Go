package dlt645

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// NewDLT645TCP97Pro
// [修改] 构造函数返回 error，不再忽略初始化/连接错误
func NewDLT645TCP97Pro(host string, port int, addr string, timeout int, fz bool) (*DLT645TCP97Pro, error) {
	dlt := new(DLT645TCP97Pro)
	dlt.Host = host
	dlt.Port = port
	if timeout <= 0 {
		timeout = 5 // [修改] 默认超时
	}
	dlt.Timeout = timeout
	dlt.Addr = addr

	// [修改] 校验表号：必须是 12 位十进制数字
	if len(addr) != 12 {
		return nil, fmt.Errorf("表号必须为12位十进制数字, got: %q", addr)
	}
	for _, c := range addr {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("表号必须全为十进制数字, got: %q", addr)
		}
	}

	address, err := dlt.strToBCD(addr)
	if err != nil {
		return nil, fmt.Errorf("表号转BCD失败: %w", err)
	}
	if fz {
		dlt.address = dlt.reverseBCD(address)
	} else {
		dlt.address = address
	}

	return dlt, nil
}

type DLT645TCP97Pro struct {
	Host    string
	Port    int
	Timeout int
	Conn    net.Conn
	Addr    string
	address []byte

	// [修改] 新增互斥锁，防止并发共用连接导致请求/响应错乱
	mu sync.Mutex
}

// CheckConn
// [修改] CheckConn 内部调用 connectLocked
func (p *DLT645TCP97Pro) CheckConn() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Conn == nil {
		return p.connectLocked()
	}
	return nil
}

func (p *DLT645TCP97Pro) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connectLocked()
}

// [修改] 内部版本：调用方必须已持锁
func (p *DLT645TCP97Pro) connectLocked() error {
	conn, err := net.DialTimeout(
		"tcp4",
		fmt.Sprintf("%s:%d", p.Host, p.Port),
		time.Duration(p.Timeout)*time.Second,
	)
	if err != nil {
		return err
	}
	p.Conn = conn
	return nil
}

func (p *DLT645TCP97Pro) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closeLocked()
}

// [修改] 内部版本：调用方必须已持锁
func (p *DLT645TCP97Pro) closeLocked() error {
	if p.Conn != nil {
		err := p.Conn.Close()
		p.Conn = nil
		return err
	}
	return nil
}

// 求和校验（不变）
func (p *DLT645TCP97Pro) calcCS(data []byte) byte {
	cs := byte(0)
	for _, b := range data {
		cs += b
	}
	return cs
}

// [修改] 修正奇数位 BCD：应该在左侧补 0，而不是把最后一个半字节放高 4 位
func (p *DLT645TCP97Pro) decimalToBCD(num int) []byte {
	str := strconv.Itoa(num)
	if len(str)%2 != 0 {
		str = "0" + str
	}
	bcd := make([]byte, len(str)/2)
	for i := 0; i < len(str); i += 2 {
		hi := str[i] - '0'
		lo := str[i+1] - '0'
		bcd[i/2] = hi<<4 | lo
	}
	return bcd
}

// 字符串转 Byte（不变）
func (p *DLT645TCP97Pro) strToBCD(data string) ([]byte, error) {
	if len(data)%2 != 0 {
		data = "0" + data
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

// 反转（不变）
func (p *DLT645TCP97Pro) reverseBCD(data []byte) []byte {
	for i, j := 0, len(data)-1; i <= j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

// [修改] 重写 decodeVal：
//  1. 修正 pos==8 死判断（原来是固定位置，改为末尾判断）
//  2. 用 ParseUint 避免 int32 溢出
//  3. 返回值改为 float64
//  4. 拷贝入参，避免修改调用方原始数据
func (p *DLT645TCP97Pro) decodeVal(rdData []byte, prec float64) (float64, error) {
	if len(rdData) == 0 {
		return 0, errors.New("数据域为空")
	}
	tmp := make([]byte, len(rdData))
	copy(tmp, rdData)

	// 反向并减 0x33
	for i, j := 0, len(tmp)-1; i <= j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j]-0x33, tmp[i]-0x33
	}

	var sb strings.Builder
	for _, b := range tmp {
		sb.WriteString(fmt.Sprintf("%02X", b))
	}
	s := sb.String()

	pos := strings.LastIndex(s, "F")
	digits := s
	if pos >= 0 {
		if pos == len(s)-1 { // [修改] 原 pos==8 是固定位置死判断
			return 0, nil
		}
		digits = s[pos+1:]
	}
	if digits == "" {
		return 0, nil
	}
	v, err := strconv.ParseUint(digits, 10, 64) // [修改] ParseUint 避免 int32 溢出
	if err != nil {
		return 0, fmt.Errorf("解析数据值失败 (%s): %w", digits, err)
	}
	return float64(v) * prec, nil // [修改] 直接返回 float64
}

// [修改] 新增：按 DLT645 帧格式读取一整帧，解决 TCP 粘包/拆包
func (p *DLT645TCP97Pro) readFrame() ([]byte, error) {
	// 设置读超时，避免永久阻塞
	deadline := time.Now().Add(time.Duration(p.Timeout) * time.Second)
	if err := p.Conn.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	var b [1]byte
	// 跳过前导 0xFE，找到帧头 0x68
	for {
		if _, err := io.ReadFull(p.Conn, b[:]); err != nil {
			return nil, fmt.Errorf("读取帧头失败: %w", err)
		}
		if b[0] == 0xFE {
			continue
		}
		if b[0] == 0x68 {
			break
		}
		return nil, fmt.Errorf("非法帧头: %02X", b[0])
	}

	frame := make([]byte, 0, 32)
	frame = append(frame, 0x68)

	// 地址 6 字节
	addr := make([]byte, 6)
	if _, err := io.ReadFull(p.Conn, addr); err != nil {
		return nil, err
	}
	frame = append(frame, addr...)

	// 第二个 0x68
	if _, err := io.ReadFull(p.Conn, b[:]); err != nil {
		return nil, err
	}
	if b[0] != 0x68 {
		return nil, fmt.Errorf("缺少第二个 0x68: %02X", b[0])
	}
	frame = append(frame, 0x68)

	// 控制码
	if _, err := io.ReadFull(p.Conn, b[:]); err != nil {
		return nil, err
	}
	frame = append(frame, b[0])

	// 数据长度
	if _, err := io.ReadFull(p.Conn, b[:]); err != nil {
		return nil, err
	}
	dataLen := int(b[0])
	frame = append(frame, b[0])

	// 数据域 + CS + 0x16
	rest := make([]byte, dataLen+2)
	if _, err := io.ReadFull(p.Conn, rest); err != nil {
		return nil, err
	}
	frame = append(frame, rest...)

	// 结束符
	if frame[len(frame)-1] != 0x16 {
		return nil, fmt.Errorf("帧结束符错误: %02X", frame[len(frame)-1])
	}

	// CS 校验：从第一个 0x68 到数据域最后一字节
	want := p.calcCS(frame[:len(frame)-2])
	got := frame[len(frame)-2]
	if want != got {
		return nil, fmt.Errorf("校验和错误: 计算=%02X, 报文=%02X", want, got)
	}
	return frame, nil
}

// [修改] 用 readFrame + 全字段校验重写读取流程
//
//	prec 与返回值均为 float64
//
// 注意：调用方必须已持有 p.mu 锁
func (p *DLT645TCP97Pro) dlt645TCP97Locked(DI string, prec float64) (ReadRes, error) {
	var res ReadRes
	// [修改] 校验 DI 格式：97 规约 DI 为 2 字节，即 4 位十六进制字符串
	if len(DI) != 4 {
		return res, fmt.Errorf("DI 必须为4位十六进制字符串(97规约), got: %q", DI)
	}
	ditR, err := p.strToBCD(DI)
	if err != nil {
		return res, fmt.Errorf("DI 转BCD失败: %w", err)
	}
	dit := p.reverseBCD(ditR)

	di := make([]byte, 0, len(dit))
	for _, d := range dit {
		di = append(di, d+0x33)
	}

	if p.Conn == nil {
		return res, errors.New("conn is nil")
	}

	// 组装请求帧
	// [修改] 97 规约请求控制码是 0x01
	req := []byte{0x68}
	req = append(req, p.address...)
	req = append(req, []byte{0x68, 0x01}...)
	req = append(req, byte(len(di)))
	req = append(req, di...)
	req = append(req, p.calcCS(req))
	req = append(req, 0x16)

	// [修改] 设置写超时
	if err := p.Conn.SetWriteDeadline(time.Now().Add(time.Duration(p.Timeout) * time.Second)); err != nil {
		return res, err
	}
	res.Write = req
	// [修改] 检查实际写入长度
	n, err := p.Conn.Write(req)
	if err != nil {
		return res, err
	}
	if n != len(req) {
		return res, fmt.Errorf("写入长度不足: %d/%d", n, len(req))
	}

	// [修改] 按帧读取（解决粘包/拆包、超时、CS 校验）
	resp, err := p.readFrame()
	if err != nil {
		return res, err
	}
	res.Read = resp

	// [修改] 校验应答地址
	if !bytes.Equal(resp[1:7], p.address) {
		return res, fmt.Errorf("应答地址不匹配: %X vs %X", resp[1:7], p.address)
	}
	// resp[7] 已由 readFrame 保证为 0x68

	ctrl := resp[8]
	dataLen := int(resp[9])
	// [修改] 长度字段与报文实际长度一致
	if 10+dataLen+2 != len(resp) {
		return res, fmt.Errorf("数据长度字段(%d)与报文长度(%d)不匹配", dataLen, len(resp))
	}
	data := resp[10 : 10+dataLen]

	// [修改] 97 规约控制码：
	//   0x81 = 正常读数据应答
	//   0xC1 = 异常应答
	switch ctrl {
	case 0x81, 0xA1:
		// 正常
	case 0xC1, 0xE1:
		if len(data) < 1 {
			return res, errors.New("异常应答无错误码")
		}
		return res, fmt.Errorf("设备异常应答, 错误码: %02X", data[0])
	default:
		return res, fmt.Errorf("未知应答控制码: %02X", ctrl)
	}

	// [修改] 校验返回的 DI 与请求 DI 一致（97 规约 DI 为 2 字节）
	if len(data) < 2 {
		return res, errors.New("应答数据域不足 2 字节 DI")
	}
	if !bytes.Equal(data[:2], di) {
		return res, fmt.Errorf("应答 DI 不匹配: %X vs %X", data[:2], di)
	}

	rdData := data[2:]
	if len(rdData) == 0 {
		return res, errors.New("应答无数据值")
	}

	// [修改] 不再吞掉 decodeVal 的 error
	val, err := p.decodeVal(rdData, prec)
	res.Val = val
	return res, err
}

// ReadTag
// [修改] ReadTag：返回 float64，加锁、管理连接、不再忽略错误
func (p *DLT645TCP97Pro) ReadTag(DI string, prec float64) (ReadRes, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Conn == nil {
		if err := p.connectLocked(); err != nil {
			return ReadRes{}, err
		}
	}
	defer p.closeLocked()

	return p.dlt645TCP97Locked(DI, prec)
}

// ReadTags
// [修改] ReadTags：输入/输出 map 均为 map[string]float64
// 加锁、收集错误、出错时重连、去掉多余 sleep
// 入:key-DI,v:prec,出:key-DI,v:值
func (p *DLT645TCP97Pro) ReadTags(DIP map[string]float64) (map[string]ReadRes, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Conn == nil {
		if err := p.connectLocked(); err != nil {
			return nil, err
		}
	}
	defer p.closeLocked()

	readList := make(map[string]ReadRes, len(DIP))
	var errs []string
	i := 0
	for di, prec := range DIP {
		read, err := p.dlt645TCP97Locked(di, prec)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", di, err))
			// [修改] 连接可能已经损坏，主动重连，避免后续全错
			_ = p.closeLocked()
			if cerr := p.connectLocked(); cerr != nil {
				return readList, fmt.Errorf("重连失败: %w (已收集错误: %s)", cerr, strings.Join(errs, "; "))
			}
		} else {
			readList[di] = read
		}
		i++
		if i < len(DIP) { // [修改] 最后一个不再 sleep
			time.Sleep(100 * time.Millisecond)
		}
	}
	if len(errs) > 0 {
		return readList, fmt.Errorf("部分读取失败: %s", strings.Join(errs, "; "))
	}
	return readList, nil
}
