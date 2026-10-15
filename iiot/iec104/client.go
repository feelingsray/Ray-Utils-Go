package iec104

import (
	"context"
	"net"
	"sync"
	"time"
)

var (
	contextTimeout    = 30 * time.Second
	dialTimeout       = 5 * time.Second
	testInterval      = 20 * time.Second
	totalCallInterval = 15 * time.Minute
	retryTimes        = 3 // 存在备用服务器时，单个服务器重试次数
)

// Client 104客户端
type IEC104Client struct {
	address    string
	subAddress string
	curAddress string
	conn       net.Conn
	cancel     context.CancelFunc
	rsn        int16
	ssn        int16
	dataChan   chan *APDU
	sendChan   chan []byte
	iFrameNum  int
	task       func(c *APDU)
	wg         *sync.WaitGroup
}

// NewClient 初始化客户端,连接失败，每隔10秒重试
func NewIEC104Client(address string, subAddress ...string) *IEC104Client {
	subAddr := ""
	if len(subAddress) == 1 && subAddress[0] != "" {
		subAddr = subAddress[0]
	}
	return &IEC104Client{
		address:    address,
		subAddress: subAddr,
		curAddress: address,
		dataChan:   make(chan *APDU, 1),
		sendChan:   make(chan []byte, 1),
		wg:         new(sync.WaitGroup),
	}
}

// Run 运行
func (c *IEC104Client) Run(parent context.Context, task func(context.Context, *APDU, string, int), company string, fre int) {
	for {
		c.conn = c.dail()
		c.sendUFrame(startDtAct)
		ctx, cancel := context.WithCancel(parent)
		c.cancel = cancel
		c.wg.Add(3)
		go c.read(ctx)
		go c.write(ctx)
		go c.handler(ctx, task, company, fre)
		c.wg.Wait()
		if c.conn != nil {
			c.conn.Close()
		}
		c.rsn = 0
		c.ssn = 0
		c.iFrameNum = 0
		time.Sleep(100 * time.Second)
	}
}

// 建立tcp连接，支持重试和主备切换
func (c *IEC104Client) dail() net.Conn {
	var conn net.Conn
	var err error
	i := -1
	for {
		conn, err = net.DialTimeout("tcp", c.curAddress, dialTimeout)
		if err != nil {
			time.Sleep(dialTimeout)
			i++
			if i == retryTimes && c.subAddress != "" {
				i = 0
				if c.curAddress == c.address {
					c.curAddress = c.subAddress
				} else {
					c.curAddress = c.address
				}
			}
		} else {
			break
		}
	}
	return conn
}

// Read 读数据
func (c *IEC104Client) read(ctx context.Context) {
	defer func() {
		c.cancel()
		c.wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.parseData()
			if err != nil {
				return
			}
		}
	}
}

// Write 写数据
func (c *IEC104Client) write(ctx context.Context) {
	defer func() {
		c.cancel()
		c.wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c.sendChan:
			_, err := c.conn.Write(data)
			if err != nil {
				return
			}
		}
	}
}

// handler 处理接收到的已解析数据
func (c *IEC104Client) handler(ctx context.Context, task func(ctx context.Context, c *APDU, p string, f int), company string, fre int) {
	defer func() {
		c.cancel()
		c.wg.Done()
	}()
	for {
		select {
		case resp := <-c.dataChan:
			go task(ctx, resp, company, fre)
		case <-ctx.Done():
			return
		}
	}
}

// ParseData 解析接收到的数据
func (c *IEC104Client) parseData() error {
	buf := make([]byte, 2)
	// 读取启动符和长度
	n, err := c.conn.Read(buf)
	if err != nil {
		return err
	}
	c.conn.SetDeadline(time.Now().Add(contextTimeout))
	if n == 0 {
		time.Sleep(10 * time.Second)
		return nil
	}
	length := int(buf[1])
	// 读取正文
	contentBuf := make([]byte, length)
	n, err = c.conn.Read(contentBuf)
	if err != nil {
		return err
	}
	// 长度不够继续读取，直至达到期望长度
	i := 1
	for n < length {
		i++
		nextLength := length - n
		nextBuf := make([]byte, nextLength)
		m, err := c.conn.Read(nextBuf)
		if err != nil {
			return err
		}
		contentBuf = append(contentBuf[:n], nextBuf[:m]...)
		n = len(contentBuf)
	}
	apdu := new(APDU)
	err = apdu.parseAPDU(contentBuf[:n])
	if err != nil {
		// c.Logger.Warnf("解析APDU异常: %v", err)
		// c.Logger.Panicln("退出程序")
		return err
	}
	switch apdu.CtrFrame.(type) {
	case IFrame:
		c.incrRsn()
		switch apdu.ASDU.TypeID {
		case MEiNA1:
			c.sendSFrame()
			c.sendTotalCall()
		case CIcNa1:
			if apdu.ASDU.Cause == 7 {
				c.sendSFrame()
			} else if apdu.ASDU.Cause == 10 {
				c.sendSFrame()
			}
		case CCiNa1:
			if apdu.ASDU.Cause == 7 {
				// c.Logger.Sound("接收电度总召唤确认帧")
			} else if apdu.ASDU.Cause == 10 {
				// c.Logger.Sound("接收电度总召唤结束帧")
			}
			c.sendSFrame()
		default:
			c.iFrameNum++
			c.dataChan <- apdu
			c.sendSFrame()
		}
	case SFrame:
	case UFrame:
		uFrame := apdu.CtrFrame.(UFrame)
		switch uFrame.cmd {
		case startDtCon:
			// c.Logger.Sound("U帧为启动确认帧，发送总召唤")
			c.sendTotalCall()
		case testFrAct:
			// c.Logger.Sound("U帧为测试激活帧,发送测试确认帧")
			c.sendUFrame(testFrCon)
		}
	default:
	}
	return nil
}

// sendUFrame 发送U帧
func (c *IEC104Client) sendUFrame(cmd [4]byte) {
	data := convertBytes(convert4BytesToSlice(cmd))
	c.sendChan <- data
}

// sendSFrame 发送S帧
func (c *IEC104Client) sendSFrame() {
	rsnBytes := parseLittleEndianUInt16(uint16(c.rsn << 1))
	sendBytes := make([]byte, 0)
	sendBytes = append(sendBytes, 0x01, 0x00)
	sendBytes = append(sendBytes, rsnBytes...)
	data := convertBytes(sendBytes)
	c.sendChan <- data
}

// sendTotalCall 发送总召唤
func (c *IEC104Client) sendTotalCall() {
	ssnBytes := parseLittleEndianUInt16(uint16(c.ssn << 1))
	rsnBytes := parseLittleEndianUInt16(uint16(c.rsn << 1))
	totalCallData := make([]byte, 0)
	totalCallData = append(totalCallData, ssnBytes...)
	totalCallData = append(totalCallData, rsnBytes...)
	totalCallData = append(totalCallData, 0x64, 0x01, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x14)
	data := convertBytes(totalCallData)
	c.sendChan <- data
}

// sendTotalCall 发送电度总召唤
func (c *IEC104Client) sendElectricityTotalCall() {
	ssnBytes := parseLittleEndianUInt16(uint16(c.ssn << 1))
	rsnBytes := parseLittleEndianUInt16(uint16(c.rsn << 1))
	totalCallData := make([]byte, 0)
	totalCallData = append(totalCallData, ssnBytes...)
	totalCallData = append(totalCallData, rsnBytes...)
	totalCallData = append(totalCallData, 0x65, 0x01, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x05)
	data := convertBytes(totalCallData)
	c.sendChan <- data
}

// incrRsn 增加rsn
func (c *IEC104Client) incrRsn() {
	c.rsn++
	if c.rsn < 0 {
		c.rsn = 0
	}
}
