package modbus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

type ByteOrder int
type DataType string

const (
	MH_AB_CD       ByteOrder = 0
	MH_CD_AB       ByteOrder = 1
	MH_BA_DC       ByteOrder = 2
	MH_DC_BA       ByteOrder = 3
	MH_AB_CD_EF_GH ByteOrder = 0
	MH_GH_EF_CD_AB ByteOrder = 1
	MH_BA_DC_FE_HG ByteOrder = 2
	MH_HG_FE_DC_BA ByteOrder = 3

	MH_INT16 DataType = "int16"
)

type ModBusHelper struct {
}

func NewModBusHelper() *ModBusHelper {
	return &ModBusHelper{}
}

/*
 * 线性缩放方法
 * @param value float64
 * offset 可能就是 rawLow
 */
func (m *ModBusHelper) LinerScaling(value interface{}, offset float64, rawHeight float64, rawLow float64, scaleHeight float64, scaleLow float64) (float64, error) {
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

/*
 * Short转换为寄存器值
 * 说明:
 *      Short类型即为int16类型,modbus中读取1个寄存器做无符号处理
 */
func (m *ModBusHelper) ShortTo(raw int16) uint16 {
	var uValue uint16 = uint16(raw)
	b := make([]byte, 2)
	if 0 == 0 {
		binary.BigEndian.PutUint16(b, uValue)
		return binary.BigEndian.Uint16(b)
	} else {
		binary.LittleEndian.PutUint16(b, uValue)
		return binary.LittleEndian.Uint16(b)
	}
}

/*
 * 寄存器值转换为short类型
 * 说明:
 *      Short类型即为int16类型,modbus中读取1个寄存器做无符号处理
 */
func (m *ModBusHelper) ToShort(raw uint16) int16 {
	return int16(raw)
}

/*
 * Int类型转换为寄存器值
 * int类型默认为int32
 */
func (m *ModBusHelper) IntTo(raw int32, order ByteOrder) ([2]uint16, error) {
	data := [2]uint16{}
	bytesBuffer := bytes.NewBuffer([]byte{})
	if order == 0 || order == 1 {
		if err := binary.Write(bytesBuffer, binary.BigEndian, raw); err != nil {
			return data, err
		} else {
			data[0] = binary.BigEndian.Uint16(bytesBuffer.Bytes()[:2])
			data[1] = binary.BigEndian.Uint16(bytesBuffer.Bytes()[2:])
		}
	} else {
		if err := binary.Write(bytesBuffer, binary.LittleEndian, raw); err != nil {
			return data, err
		} else {
			data[0] = binary.LittleEndian.Uint16(bytesBuffer.Bytes()[:2])
			data[1] = binary.LittleEndian.Uint16(bytesBuffer.Bytes()[2:])
		}
	}
	if order == 1 || order == 3 {
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}
	return data, nil
}

/*
 * 寄存器值转换为int类型
 * int类型默认为int32
 */
func (m *ModBusHelper) ToInt(raw [2]uint16, order ByteOrder) (int32, error) {
	var data int32
	if order == 1 || order == 3 {
		for i, j := 0, len(raw)-1; i < j; i, j = i+1, j-1 {
			raw[i], raw[j] = raw[j], raw[i]
		}
	}
	lb := make([]byte, 2)
	hb := make([]byte, 2)
	if order == 0 || order == 1 {
		binary.BigEndian.PutUint16(lb, raw[0])
		binary.BigEndian.PutUint16(hb, raw[1])
		lb = append(lb, hb...)
		if err := binary.Read(bytes.NewBuffer(lb), binary.BigEndian, &data); err != nil {
			return data, err
		}
	} else {
		binary.LittleEndian.PutUint16(lb, raw[0])
		binary.LittleEndian.PutUint16(hb, raw[1])
		lb = append(lb, hb...)
		if err := binary.Read(bytes.NewBuffer(lb), binary.LittleEndian, &data); err != nil {
			return data, err
		}
	}
	return data, nil
}

/*
 * 不一定对哈
 * Int类型转换为寄存器值
 * int类型默认为int32
 */
func (m *ModBusHelper) FloatTo(raw float32, order ByteOrder) ([2]uint16, error) {
	data := [2]uint16{}
	bytesBuffer := bytes.NewBuffer([]byte{})
	if order == 0 || order == 1 {
		if err := binary.Write(bytesBuffer, binary.BigEndian, raw); err != nil {
			return data, err
		} else {
			data[0] = binary.BigEndian.Uint16(bytesBuffer.Bytes()[:2])
			data[1] = binary.BigEndian.Uint16(bytesBuffer.Bytes()[2:])
		}
	} else {
		if err := binary.Write(bytesBuffer, binary.LittleEndian, raw); err != nil {
			return data, err
		} else {
			data[0] = binary.LittleEndian.Uint16(bytesBuffer.Bytes()[:2])
			data[1] = binary.LittleEndian.Uint16(bytesBuffer.Bytes()[2:])
		}
	}
	if order == 1 || order == 3 {
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}
	return data, nil
}

/*
 * 寄存器值转换为int类型
 * int类型默认为int32
 */
func (m *ModBusHelper) ToFloat(raw [2]uint16, order ByteOrder) (float32, error) {
	var data float32
	if order == 1 || order == 3 {
		for i, j := 0, len(raw)-1; i < j; i, j = i+1, j-1 {
			raw[i], raw[j] = raw[j], raw[i]
		}
	}
	lb := make([]byte, 2)
	hb := make([]byte, 2)
	if order == 0 || order == 1 {
		binary.BigEndian.PutUint16(lb, raw[0])
		binary.BigEndian.PutUint16(hb, raw[1])
		lb = append(lb, hb...)
		if err := binary.Read(bytes.NewBuffer(lb), binary.BigEndian, &data); err != nil {
			return data, err
		}
	} else {
		binary.LittleEndian.PutUint16(lb, raw[0])
		binary.LittleEndian.PutUint16(hb, raw[1])
		lb = append(lb, hb...)
		if err := binary.Read(bytes.NewBuffer(lb), binary.LittleEndian, &data); err != nil {
			return data, err
		}
	}
	return data, nil
}
