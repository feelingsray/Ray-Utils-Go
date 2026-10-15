package s7

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	"github.com/robinson/gos7"
)

const (
	// Area ID
	s7areape = 0x81 // process inputs
	s7areapa = 0x82 // process outputs
	s7areamk = 0x83 // Merkers
	s7areadb = 0x84 // DB
	s7areact = 0x1C // counters
	s7areatm = 0x1D // timers

	// Word Length
	s7wlbit     = 0x01 // Bit (inside a word)
	s7wlbyte    = 0x02 // Byte (8 bit)
	s7wlChar    = 0x03
	s7wlword    = 0x04 // Word (16 bit)
	s7wlint     = 0x05
	s7wldword   = 0x06 // Double Word (32 bit)
	s7wldint    = 0x07
	s7wlreal    = 0x08 // Real (32 bit float)
	s7wlcounter = 0x1C // Counter (16 bit)
	s7wltimer   = 0x1D // Timer (16 bit)

	// PLC Status
	s7CpuStatusUnknown = 0
	s7CpuStatusRun     = 8
	s7CpuStatusStop    = 4

	// size header
	sizeHeaderRead  int = 31 // Header Size when Reading
	sizeHeaderWrite int = 35 // Header Size when Writing

	// Result transport size
	tsResBit   = 3
	tsResByte  = 4
	tsResInt   = 5
	tsResReal  = 7
	tsResOctet = 9
)

type MyGos7Client struct {
	packager    gos7.Packager
	transporter gos7.Transporter
}

func NewMyGos7Client(handler gos7.ClientHandler) *MyGos7Client {
	return &MyGos7Client{packager: handler, transporter: handler}
}

func (mb *MyGos7Client) Read(variable string, buffer []byte) (value any, err error) {
	variable = strings.ToUpper(variable)             // upper
	variable = strings.ReplaceAll(variable, " ", "") // remove spaces
	if variable == "" {
		err = errors.New("input variable is empty, variable should be S7 syntax")
		return
	}
	// var area, dbNumber, start, amount, wordLen int
	switch valueArea := variable[0:2]; valueArea {
	case "VB": // s7200 smart 使用地址

	case "VW": // s7200 smart 使用地址
	case "VD": // s7200 smart 使用地址
		startByte, _ := strconv.ParseInt(variable[2:], 10, 16)
		err = mb.AGReadDB(1, int(startByte), 4, buffer)
		if err != nil {
			return
		}
		// helper := gos7.Helper{}
		// helper.GetValueAt(buffer, 0, value)
		return
	case "EB": // input byte
	case "EW": // input word
	case "ED": // Input double-word
	case "AB": // Output byte
	case "AW": // Output word
	case "AD": // Output double-word
	case "MB": // Memory byte
	case "MW": // Memory word
	case "MD": // Memory double-word
	case "DB": // Data Block
		dbArray := strings.Split(variable, ".")
		if len(dbArray) < 2 {
			err = errors.New("Db Area read variable should not be empty")
			return
		}
		dbNo, _ := strconv.ParseInt(dbArray[0][2:], 10, 16)
		dbIndex, _ := strconv.ParseInt(dbArray[1][3:], 10, 16)
		dbType := dbArray[1][0:3]

		switch dbType {
		case "DBB": // byte
			err = mb.AGReadDB(int(dbNo), int(dbIndex), 1, buffer)
			value = buffer[0]
			return
		case "DBW": // word
			err = mb.AGReadDB(int(dbNo), int(dbIndex), 2, buffer)
			value = binary.BigEndian.Uint16(buffer[0:])
			return
		case "DBD": // dword
			err = mb.AGReadDB(int(dbNo), int(dbIndex), 4, buffer)
			value = binary.BigEndian.Uint32(buffer[0:])
			return
		case "DBX": // bit
			mBit, _ := strconv.ParseInt(dbArray[2][0:], 10, 16)
			if mBit > 7 || mBit < 0 {
				err = errors.New("Db read bit is invalid")
				return
			}
			err = mb.AGReadDB(int(dbNo), int(dbIndex), 1, buffer)
			// mask := []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
			// value = buffer[0] & mask[mBit]
			return
		default:
			err = errors.New("error when parsing dbtype")
			return
		}
	default:
		switch otherArea := variable[0:1]; otherArea {
		case "E":
		case "I": // input
			mArray := strings.Split(variable, ".")

			startByte, _ := strconv.ParseInt(mArray[0][1:], 10, 16)
			err = mb.AGReadEB(int(startByte), 1, buffer)
			if err != nil {
				return
			}
		case "A":
		case "0": // output
		case "Q": // output
			mArray := strings.Split(variable, ".")
			startByte, _ := strconv.ParseInt(mArray[0][1:], 10, 16)
			err = mb.AGReadAB(int(startByte), 1, buffer)
			if err != nil {
				return
			}
		case "M": // memory
			mArray := strings.Split(variable, ".")

			startByte, _ := strconv.ParseInt(mArray[0][1:], 10, 16)
			err = mb.AGReadMB(int(startByte), 1, buffer)
			if err != nil {
				return
			}
			// helper := gos7.Helper{}
			// helper.GetValueAt(buffer, 0, value)
			return
		case "T": // timer
			startByte, _ := strconv.ParseInt(variable[1:], 10, 16)
			err = mb.AGReadTM(int(startByte), 1, buffer)
			if err != nil {
				return
			}
			helper := gos7.Helper{}
			helper.GetValueAt(buffer, 0, value)
			return
		case "Z":
		case "C": // counter
			startByte, _ := strconv.ParseInt(variable[1:], 10, 16)
			err = mb.AGReadCT(int(startByte), 1, buffer)
			if err != nil {
				return
			}
			helper := gos7.Helper{}
			helper.GetValueAt(buffer, 0, value)
			return
		case "V":
			startByte, _ := strconv.ParseInt(variable[1:], 10, 16)
			err = mb.AGReadDB(1, int(startByte), 1, buffer)
			if err != nil {
				return
			}
			helper := gos7.Helper{}
			helper.GetValueAt(buffer, 0, value)
			return
		default:
			err = errors.New("error when parsing db area")
			return
		}
	}
	return
	// }
}

// implement of the interface AGReadDB
func (mb *MyGos7Client) AGReadDB(dbnumber int, start int, size int, buffer []byte) (err error) {
	return mb.readArea(s7areadb, dbnumber, start, size, s7wlbyte, buffer)
}

// implement of the interface AGReadDB
func (mb *MyGos7Client) AGReadPE(dbnumber int, start int, size int, buffer []byte) (err error) {
	return mb.readArea(s7areape, dbnumber, start, size, s7wlbyte, buffer)
}

// implement of the interface AGWriteDB
func (mb *MyGos7Client) AGWriteDB(dbNumber int, start int, size int, buffer []byte) (err error) {
	return mb.writeArea(s7areadb, dbNumber, start, size, s7wlbyte, buffer)
}

// implement of the interface AGReadMB
func (mb *MyGos7Client) AGReadMB(start int, size int, buffer []byte) (err error) {
	return mb.readArea(s7areamk, 0, start, size, s7wlbyte, buffer)
}

// implement of the interface AGWriteMB
func (mb *MyGos7Client) AGWriteMB(start int, size int, buffer []byte) (err error) {
	return mb.writeArea(s7areamk, 0, start, size, s7wlbyte, buffer)
}

// implement of the interface AGReadEB
func (mb *MyGos7Client) AGReadEB(start int, size int, buffer []byte) (err error) {
	return mb.readArea(s7areape, 0, start, size, s7wlbyte, buffer)
}

// implement of the interface AGWriteEB
func (mb *MyGos7Client) AGWriteEB(start int, size int, buffer []byte) (err error) {
	return mb.writeArea(s7areape, 0, start, size, s7wlbyte, buffer)
}

// implement of the interface AGReadAB
func (mb *MyGos7Client) AGReadAB(start int, size int, buffer []byte) (err error) {
	return mb.readArea(s7areapa, 0, start, size, s7wlbyte, buffer)
}

// implement of the interface AGWriteAB
func (mb *MyGos7Client) AGWriteAB(start int, size int, buffer []byte) (err error) {
	return mb.writeArea(s7areapa, 0, start, size, s7wlbyte, buffer)
}

// implement of the interface AGReadTM - read timer
func (mb *MyGos7Client) AGReadTM(start int, amount int, buffer []byte) (err error) {
	sbuffer := make([]byte, amount*2)
	err = mb.readArea(s7areatm, 0, start, amount, s7wltimer, sbuffer)
	if err == nil {
		for c := range amount {
			buffer[c] = byte(uint16(sbuffer[c*2+1])<<8 + uint16(sbuffer[c*2]))
		}
	}
	return err
}

// implement of the interface AGWriteTM - write timer
func (mb *MyGos7Client) AGWriteTM(start int, amount int, buffer []byte) (err error) {
	sbuffer := make([]byte, amount*2)
	for c := range amount {
		sbuffer[c*2+1] = byte((uint(buffer[c]) & uint(0xFF00)) >> 8)
		sbuffer[c*2] = buffer[c] & 0x00FF
	}
	err = mb.writeArea(s7areatm, 0, start, amount, s7wltimer, sbuffer)
	return err
}

// implement of the interface AGReadCT - read counter
func (mb *MyGos7Client) AGReadCT(start int, amount int, buffer []byte) (err error) {
	sbuffer := make([]byte, amount*2)
	err = mb.readArea(s7areact, 0, start, amount, s7wlcounter, sbuffer)
	if err == nil {
		for c := range amount {
			buffer[c] = byte(uint(sbuffer[c*2+1])<<8 + uint(sbuffer[c*2]))
		}
	}
	return err
}

// implement of the interface AGWriteCT - write counter
func (mb *MyGos7Client) AGWriteCT(start int, amount int, buffer []byte) (err error) {
	sbuffer := make([]byte, amount*2)
	for c := range amount {
		sbuffer[c*2+1] = byte((uint(buffer[c]) & uint(0xFF00)) >> 8)
		sbuffer[c*2] = buffer[c] & 0x00FF
	}
	err = mb.writeArea(s7areact, 0, start, amount, s7wlcounter, sbuffer)
	return err
}

func (mb *MyGos7Client) readArea(area int, dbNumber int, start int, amount int, wordLen int, buffer []byte) (err error) {
	var address, numElements, maxElements, totElements, sizeRequested int
	offset := 0
	wordSize := 1
	// Some adjustment
	if area == s7areact {
		wordLen = s7wlcounter
	}
	if area == s7areatm {
		wordLen = s7wltimer
	}
	// Calc Word size
	wordSize = dataSizeByte(wordLen)
	if wordSize == 0 {
		return errors.New(gos7.ErrorText(errIsoInvalidDataSize))
	}

	if wordLen == s7wlbit {
		amount = 1 // Only 1 bit can be transferred at time
	} else {
		if wordLen != s7wlcounter && wordLen != s7wltimer {
			amount = amount * wordSize
			wordSize = 1
			wordLen = s7wlbyte
		}
	}

	tt, _ := any(mb.transporter).(*gos7.TCPClientHandler)

	maxElements = (tt.PDULength - 18) / wordSize // 18 = Reply telegram header //lth note here
	totElements = amount
	for totElements > 0 && err == nil {
		numElements = min(totElements, maxElements)

		sizeRequested = numElements * wordSize
		// Agent the telegram
		requestData := make([]byte, sizeHeaderRead)
		copy(requestData[0:], s7ReadWriteTelegram[0:])
		request := gos7.NewProtocolDataUnit(requestData)
		// Set DB Number
		request.Data[27] = byte(area)
		// Set Area
		if area == s7areadb {
			binary.BigEndian.PutUint16(request.Data[25:], uint16(dbNumber))
			// SetWordAt(request.Data, 25, uint16(DBNumber))
		}

		// Adjusts Start and word length
		if wordLen == s7wlbit || wordLen == s7wlcounter || wordLen == s7wltimer {
			address = start
			request.Data[22] = byte(wordLen)
		} else {
			address = start << 3
		}
		// Num elements
		binary.BigEndian.PutUint16(request.Data[23:], uint16(numElements))
		// SetWordAt(request.Data, 23, uint16(numElements))
		// AddressLocal into the PLC (only 3 bytes)
		request.Data[30] = byte(address & 0x0FF)
		address = address >> 8
		request.Data[29] = byte(address & 0x0FF)
		address = address >> 8
		request.Data[28] = byte(address & 0x0FF)
		var response *gos7.ProtocolDataUnit
		response, sendError := mb.send(&request)
		err = sendError

		if err == nil {
			if size := len(response.Data); size < 25 {
				err = errors.New(gos7.ErrorText(errIsoInvalidDataSize))
			} else {
				if response.Data[21] != 0xFF {
					err = errors.New(gos7.ErrorText(CPUError(uint(response.Data[21]))))
				} else {
					// copy response to buffer
					copy(buffer[offset:offset+sizeRequested], response.Data[25:25+sizeRequested])
					offset += sizeRequested
				}
			}
		}
		totElements -= numElements
		start += numElements * wordSize
	}
	return
}

// writeArea write generic area into PLC with following parameters:
// 1.area: s7areape/s7areapa/s7areamk/s7areadb/s7areact/s7areatm
// 2.dbnumber: specify dbnumber, to use in write DB area, otherwise = 0
// 3.start: start of the address
// 4.amount: amount of the address
// 5.wordlen: bit/byte/word/dword/real/counter/timer
// 6.buffer: a byte array input for writing
func (mb *MyGos7Client) writeArea(area int, dbnumber int, start int, amount int, wordlen int, buffer []byte) (err error) {
	var address, numElements, maxElements, totElements, dataSize, isoSize, length int
	offset := 0
	wordSize := 1

	// Some adjustment
	if area == s7areact {
		wordlen = s7wlcounter
	}
	if area == s7areatm {
		wordlen = s7wltimer
	}

	// Calc Word size
	wordSize = dataSizeByte(wordlen)
	if wordSize == 0 {
		return errors.New(gos7.ErrorText(errIsoInvalidDataSize))
	}

	if wordlen == s7wlbit {
		amount = 1 // Only 1 bit can be transferred at time
	} else {
		if wordlen != s7wlcounter && wordlen != s7wltimer {
			amount = amount * wordSize
			wordSize = 1
			wordlen = s7wlbyte
		}
	}
	tt, _ := any(mb.transporter).(*gos7.TCPClientHandler)
	maxElements = (tt.PDULength - 35) / wordSize // 35 = Reply telegram header
	totElements = amount
	for totElements > 0 && err == nil {
		numElements = min(totElements, maxElements)
		dataSize = numElements * wordSize
		isoSize = sizeHeaderWrite + dataSize

		// Agent the telegram
		requestData := make([]byte, sizeHeaderWrite)
		copy(requestData[0:], s7ReadWriteTelegram[0:])

		request := gos7.NewProtocolDataUnit(requestData)
		// Whole telegram Size
		binary.BigEndian.PutUint16(request.Data[2:], uint16(isoSize))
		// SetWordAt(request.Data, 2, uint16(isoSize))
		// Data length
		length = dataSize + 4
		binary.BigEndian.PutUint16(request.Data[15:], uint16(length))
		// SetWordAt(request.Data, 15, uint16(length))
		// Function
		request.Data[17] = byte(0x05)
		// Set DB Number
		request.Data[27] = byte(area)
		if area == s7areadb {
			binary.BigEndian.PutUint16(request.Data[25:], uint16(dbnumber))
			// SetWordAt(request.Data, 25, uint16(dbnumber))
		}
		// Adjusts start and word length
		if wordlen == s7wlbit || wordlen == s7wlcounter || wordlen == s7wltimer {
			address = start
			length = dataSize
			request.Data[22] = byte(wordlen)
		} else {
			address = start << 3
			length = dataSize << 3
		}

		// Num elements
		binary.BigEndian.PutUint16(request.Data[23:], uint16(numElements))
		// SetWordAt(request.Data, 23, uint16(numElements))
		// address into the PLC
		request.Data[30] = byte(address & 0x0FF)
		address = address >> 8
		request.Data[29] = byte(address & 0x0FF)
		address = address >> 8
		request.Data[28] = byte(address & 0x0FF)

		// Transport Size
		switch wordlen {
		case s7wlbit:
			request.Data[32] = tsResBit
		case s7wlcounter:
		case s7wltimer:
			request.Data[32] = tsResOctet
		default:
			request.Data[32] = tsResByte // byte/word/dword etc.
		}
		// length
		// SetWordAt(request.Data, 33, uint16(length))
		binary.BigEndian.PutUint16(request.Data[33:], uint16(length))

		// expand values into array
		request.Data = append(request.Data[:35], append(buffer[offset:offset+dataSize], request.Data[35:]...)...)
		response, sendError := mb.send(&request)
		err = sendError
		if err == nil {
			if length = len(response.Data); length == 22 {
				if response.Data[21] != byte(0xFF) {
					err = errors.New(gos7.ErrorText(CPUError(uint(response.Data[21]))))
				}
			} else {
				err = errors.New(gos7.ErrorText(errIsoInvalidPDU))
			}
		}
		offset += dataSize
		totElements -= numElements
		start += numElements * wordSize
	}
	return
}

// implement of the interface AGReadDB

func (mb *MyGos7Client) send(request *gos7.ProtocolDataUnit) (response *gos7.ProtocolDataUnit, err error) {
	dataResponse, err := mb.transporter.Send(request.Data)
	if err != nil {
		return
	}

	if err = mb.packager.Verify(request.Data, dataResponse); err != nil {
		return
	}
	if len(dataResponse) == 0 {
		// Empty response
		err = errors.New("s7: response data is empty")
		return
	}
	response = &gos7.ProtocolDataUnit{
		Data: dataResponse,
	}
	// check for error if any
	err = responseError(response)
	return response, err
}

// responseError get response error from pdu return S7Error with high and low byte
func responseError(response *gos7.ProtocolDataUnit) error {
	s7Error := &gos7.S7Error{}
	if len(response.Data) > 0 {
		switch int(response.Data[1]) {
		case 1:
		case 7:
			s7Error.High = response.Data[2]
			s7Error.Low = response.Data[3]
			break
		case 2:
		case 3:
			s7Error.High = response.Data[10]
			s7Error.Low = response.Data[11]
			break
		default:
			return nil
		}
	}
	return s7Error
}

// read generic area, pass result into a buffer
// dataSize to number of byte accordingly
func dataSizeByte(wordLength int) int {
	switch wordLength {
	case s7wlbit:
		return 1
	case s7wlbyte:
		return 1
	case s7wlChar:
		return 1
	case s7wlword:
		return 2
	case s7wlint:
		return 2
	case s7wlcounter:
		return 2
	case s7wltimer:
		return 2
	case s7wldword:
		return 4
	case s7wldint:
		return 4
	case s7wlreal:
		return 4
	default:
		return 0
	}
}
