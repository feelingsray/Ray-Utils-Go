package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
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

func (m *ModBusHelper) ByteTo(data []byte, dataType string, displayType string) (interface{}, error) {

	if dataType == "float" {
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
	} else if dataType == "long" {
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
	} else if dataType == "double" {
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
	} else {
		if len(data) != 2 {
			return nil, errors.New("数据长度不正确,应为2个8位")
		}
		if displayType == "AB" {
			bits := binary.BigEndian.Uint16(data)
			return int16(bits), nil
		} else {
			bits := binary.LittleEndian.Uint16(data)
			return int16(bits), nil
		}
	}
	return nil, nil
}

func (m *ModBusHelper) ToByte(data interface{}, dataType string, displayType string) ([]byte, error) {

	if dataType == "float" {
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
	} else if dataType == "long" {
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
			fmt.Println([]byte{bytes[2], bytes[3], bytes[0], bytes[1]})
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
	} else if dataType == "double" {
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
	} else {
		if reflect.TypeOf(data).Kind() != reflect.Int16 {
			return nil, errors.New("数据类型不匹配")
		}
		if displayType == "AB" {
			bytes := make([]byte, 2)
			binary.BigEndian.PutUint16(bytes, uint16(data.(int16)))
			return bytes, nil
		} else {
			bytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(bytes, uint16(data.(int16)))
			return bytes, nil
		}
	}
	return nil, nil
}
