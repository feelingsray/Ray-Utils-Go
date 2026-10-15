package s7

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/robinson/gos7"
)

type S7Client struct {
	handler     *gos7.TCPClientHandler
	client      *MyGos7Client
	ipAddr      string
	rack        int
	slot        int
	timeout     int
	isConnected bool
}

func NewS7200SmartClient(ip string, port string, rack int, slot int, timeout int) *S7Client {
	ipAddr := ip
	if port != "" {
		ipAddr += ":" + port
	}
	return &S7Client{
		handler:     nil,
		client:      nil,
		timeout:     timeout,
		ipAddr:      ipAddr,
		rack:        rack,
		slot:        slot,
		isConnected: false,
	}
}

func (s *S7Client) Connect() error {
	if s.handler == nil {
		s.handler = gos7.NewTCPClientHandler(s.ipAddr, s.rack, s.slot)
		s.handler.Timeout = time.Duration(s.timeout) * time.Second
		s.handler.IdleTimeout = time.Duration(s.timeout) * time.Second
		// s.handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	}
	if s.client == nil {
		s.client = NewMyGos7Client(s.handler)
	}
	err := s.handler.Connect()
	if err != nil {
		s.isConnected = false
		return err
	}
	s.isConnected = true
	return nil
}

func (s *S7Client) close() error {
	s.isConnected = false
	return s.handler.Close()
}

func (s *S7Client) ReConnect() error {
	if s.isConnected {
		return nil
	}
	err := s.Connect()
	if err != nil {
		s.isConnected = false
	}
	return err
}

func (s *S7Client) ReadDB(dbNumber int, start int, size int) ([]byte, error) {
	err := s.ReConnect()
	if err != nil {
		return []byte{}, err
	}
	buffer := make([]byte, size)

	err = s.client.AGReadDB(dbNumber, start, size, buffer)
	if err != nil {
		s.close()
		return []byte{}, err
	}
	return buffer, err
}

func (s *S7Client) ReadV(start int, size int) ([]byte, error) {
	err := s.ReConnect()
	if err != nil {
		return []byte{}, err
	}
	buffer := make([]byte, size)
	err = s.client.AGReadDB(1, start, size, buffer)
	if err != nil {
		s.close()
		return []byte{}, err
	}
	return buffer, err
}

func (s *S7Client) ReadI(start int, size int) ([]byte, error) {
	err := s.ReConnect()
	if err != nil {
		return []byte{}, err
	}
	buffer := make([]byte, size)
	err = s.client.AGReadEB(start, size, buffer)
	if err != nil {
		s.close()
		return []byte{}, err
	}
	return buffer, err
}

func (s *S7Client) ReadQ(start int, size int) ([]byte, error) {
	err := s.ReConnect()
	if err != nil {
		return []byte{}, err
	}
	buffer := make([]byte, size)
	err = s.client.AGReadAB(start, size, buffer)
	if err != nil {
		s.close()
		return []byte{}, err
	}
	return buffer, err
}

func (s *S7Client) ReadM(start int, size int) ([]byte, error) {
	err := s.ReConnect()
	if err != nil {
		return []byte{}, err
	}
	buffer := make([]byte, size)
	err = s.client.AGReadMB(start, size, buffer)
	if err != nil {
		s.close()
		return []byte{}, err
	}
	return buffer, err
}

func (s *S7Client) WriteByAddress(address string, value any) error {
	listAddress := strings.Split(address[1:len(address)-1], ",")
	listAddress[2] = strings.ToUpper(listAddress[2])             // upper
	listAddress[2] = strings.ReplaceAll(listAddress[2], " ", "") // remove spaces
	if listAddress[2][:2] == "DB" {
		dbArray := strings.Split(listAddress[2], ".")
		dbNum, _ := strconv.Atoi(dbArray[0][2:])
		if len(dbArray) == 2 {
			if dbArray[1][2] == 'B' {
				start, err := strconv.Atoi(dbArray[1][3:])
				if err != nil {
					return err
				}
				writeByte := Int8ToByte(value.(int8))
				return s.WriteDB(dbNum, start, 1, writeByte)
			} else if dbArray[1][2] == 'W' {
				start, err := strconv.Atoi(dbArray[1][3:])
				if err != nil {
					return err
				}
				var writeByte []byte
				if listAddress[0][:3] == "int" {
					writeByte = Int16ToByte(value.(int16))
				} else if listAddress[0][:4] == "uint" {
					writeByte = Uint16ToByte(value.(uint16))
				} else {
					return errors.New("S7数据类型错误，请重新检查地址")
				}
				// writeByte := Int16ToByte(value.(int16))
				return s.WriteDB(dbNum, start, 2, writeByte)
			} else if dbArray[1][2] == 'D' {
				start, err := strconv.Atoi(dbArray[1][3:])
				if err != nil {
					return err
				}
				var writeByte []byte
				if listAddress[0][:3] == "int" {
					writeByte = Int32ToByte(value.(int32), false)
				} else if listAddress[0][:4] == "uint" {
					writeByte = Uint32ToByte(value.(uint32), false)
				} else if listAddress[0][:5] == "float" {
					writeByte = Float32ToByte(value.(float32), false)
				} else {
					return errors.New("S7数据类型错误，请重新检查地址")
				}
				return s.WriteDB(dbNum, start, 4, writeByte)
			} else {
				return errors.New("S7地址错误，请重新检查地址")
			}
		} else if len(dbArray) == 3 {
			if dbArray[1][2] == 'X' {
				start, err := strconv.Atoi(dbArray[1][3:])
				if err != nil {
					return err
				}
				bit, err := strconv.Atoi(dbArray[2])
				if err != nil {
					return err
				}
				buf, err := s.ReadDB(dbNum, start, 1)
				if err != nil {
					return err
				}
				int16Value := ByteToInt16([]byte{0, buf[0]})
				writeInt16Value := EditBitInt16(int16Value, int16(bit), int16(value.(int8)))
				writeByte := Int8ToByte(int8(writeInt16Value))
				return s.WriteDB(dbNum, start, 1, writeByte)
			}
			return errors.New("s7地址错误，请重新检查地址")
		} else {
			return errors.New("s7地址错误，请重新检查地址")
		}
	} else if listAddress[2][:1] == "V" {
		if listAddress[2][1] == 'B' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			writeByte := Int8ToByte(value.(int8))
			return s.WriteV(start, 1, writeByte)
		} else if listAddress[2][1] == 'W' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int16ToByte(value.(int16))
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint16ToByte(value.(uint16))
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteV(start, 2, writeByte)
		} else if listAddress[2][1] == 'D' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int32ToByte(value.(int32), false)
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint32ToByte(value.(uint32), false)
			} else if listAddress[0][:5] == "float" {
				writeByte = Float32ToByte(value.(float32), false)
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteV(start, 4, writeByte)
		} else {
			vArray := strings.Split(listAddress[2], ".")
			if len(vArray) == 2 {
				start, err := strconv.Atoi(vArray[0][1:])
				if err != nil {
					return err
				}
				bit, err := strconv.Atoi(vArray[1])
				if err != nil {
					return err
				}
				buf, err := s.ReadV(start, 1)
				if err != nil {
					return err
				}
				int16Value := ByteToInt16([]byte{0, buf[0]})
				writeInt16Value := EditBitInt16(int16Value, int16(bit), int16(value.(int8)))
				writeByte := Int8ToByte(int8(writeInt16Value))
				return s.WriteV(start, 1, writeByte)
			}
			return errors.New("s7地址错误，请重新检查地址")
		}
	} else if listAddress[2][:1] == "Q" {
		if listAddress[2][1] == 'B' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			writeByte := Int8ToByte(value.(int8))
			return s.WriteQ(start, 1, writeByte)
		} else if listAddress[2][1] == 'W' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int16ToByte(value.(int16))
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint16ToByte(value.(uint16))
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteQ(start, 2, writeByte)
		} else if listAddress[2][1] == 'D' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int32ToByte(value.(int32), false)
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint32ToByte(value.(uint32), false)
			} else if listAddress[0][:5] == "float" {
				writeByte = Float32ToByte(value.(float32), false)
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteQ(start, 4, writeByte)
		} else {
			vArray := strings.Split(listAddress[2], ".")
			if len(vArray) == 2 {
				start, err := strconv.Atoi(vArray[0][1:])
				if err != nil {
					return err
				}
				bit, err := strconv.Atoi(vArray[1])
				if err != nil {
					return err
				}
				buf, err := s.ReadQ(start, 1)
				if err != nil {
					return err
				}
				int16Value := ByteToInt16([]byte{0, buf[0]})
				writeInt16Value := EditBitInt16(int16Value, int16(bit), int16(value.(int8)))
				writeByte := Int8ToByte(int8(writeInt16Value))
				return s.WriteQ(start, 1, writeByte)
			}
			return errors.New("s7地址错误，请重新检查地址")
		}
	} else if listAddress[2][:1] == "M" {
		if listAddress[2][1] == 'B' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			writeByte := Int8ToByte(value.(int8))
			return s.WriteM(start, 1, writeByte)
		} else if listAddress[2][1] == 'W' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int16ToByte(value.(int16))
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint16ToByte(value.(uint16))
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteM(start, 2, writeByte)
		} else if listAddress[2][1] == 'D' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int32ToByte(value.(int32), false)
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint32ToByte(value.(uint32), false)
			} else if listAddress[0][:5] == "float" {
				writeByte = Float32ToByte(value.(float32), false)
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteM(start, 4, writeByte)
		} else {
			vArray := strings.Split(listAddress[2], ".")
			if len(vArray) == 2 {
				start, err := strconv.Atoi(vArray[0][1:])
				if err != nil {
					return err
				}
				bit, err := strconv.Atoi(vArray[1])
				if err != nil {
					return err
				}
				buf, err := s.ReadM(start, 1)
				if err != nil {
					return err
				}
				int16Value := ByteToInt16([]byte{0, buf[0]})
				writeInt16Value := EditBitInt16(int16Value, int16(bit), int16(value.(int8)))
				writeByte := Int8ToByte(int8(writeInt16Value))
				return s.WriteM(start, 1, writeByte)
			}
			return errors.New("s7地址错误，请重新检查地址")
		}
	} else if listAddress[2][:1] == "I" {
		if listAddress[2][1] == 'B' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			writeByte := Int8ToByte(value.(int8))
			return s.WriteI(start, 1, writeByte)
		} else if listAddress[2][1] == 'W' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int16ToByte(value.(int16))
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint16ToByte(value.(uint16))
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteI(start, 2, writeByte)
		} else if listAddress[2][1] == 'D' {
			start, err := strconv.Atoi(listAddress[2][2:])
			if err != nil {
				return err
			}
			var writeByte []byte
			if listAddress[0][:3] == "int" {
				writeByte = Int32ToByte(value.(int32), false)
			} else if listAddress[0][:4] == "uint" {
				writeByte = Uint32ToByte(value.(uint32), false)
			} else if listAddress[0][:5] == "float" {
				writeByte = Float32ToByte(value.(float32), false)
			} else {
				return errors.New("s7数据类型错误，请重新检查地址")
			}
			return s.WriteI(start, 4, writeByte)
		} else {
			vArray := strings.Split(listAddress[2], ".")
			if len(vArray) == 2 {
				start, err := strconv.Atoi(vArray[0][1:])
				if err != nil {
					return err
				}
				bit, err := strconv.Atoi(vArray[1])
				if err != nil {
					return err
				}
				buf, err := s.ReadI(start, 1)
				if err != nil {
					return err
				}
				int16Value := ByteToInt16([]byte{0, buf[0]})
				writeInt16Value := EditBitInt16(int16Value, int16(bit), int16(value.(int8)))
				writeByte := Int8ToByte(int8(writeInt16Value))
				return s.WriteI(start, 1, writeByte)
			}
			return errors.New("s7地址错误，请重新检查地址")
		}
	} else {
		return errors.New("s7地址错误，请重新检查地址")
	}
}

func (s *S7Client) WriteDB(dbNumber int, start int, size int, buffer []byte) error {
	err := s.ReConnect()
	if err != nil {
		return err
	}
	return s.client.AGWriteDB(dbNumber, start, size, buffer)
}

func (s *S7Client) WriteV(start int, size int, buffer []byte) error {
	err := s.ReConnect()
	if err != nil {
		return err
	}
	return s.client.AGWriteDB(1, start, size, buffer)
}

func (s *S7Client) WriteI(start int, size int, buffer []byte) error {
	err := s.ReConnect()
	if err != nil {
		return err
	}
	return s.client.AGWriteEB(start, size, buffer)
}

func (s *S7Client) WriteQ(start int, size int, buffer []byte) error {
	err := s.ReConnect()
	if err != nil {
		return err
	}
	return s.client.AGWriteAB(start, size, buffer)
}

func (s *S7Client) WriteM(start int, size int, buffer []byte) error {
	err := s.ReConnect()
	if err != nil {
		return err
	}
	return s.client.AGWriteMB(start, size, buffer)
}

func (s *S7Client) ReadAddress(address string) error {
	err := s.ReConnect()
	if err != nil {
		return err
	}
	// buf := make([]byte, 255)
	// err = s.client.AGReadDB(11, 120, 4, buf)
	// log.Println(buf)
	// status, err := s.client.PLCGetStatus()
	// log.Println(status)
	// info, err := s.client.GetCPInfo()
	// log.Println(info, err)
	buf := make([]byte, 255)
	_, err = s.client.Read(address, buf)
	return err
}
