package s7

import (
	"reflect"
	"strconv"
	"time"
)

type DataValue struct {
	Time         string  `json:"time"`
	Value        any     `json:"value"`
	Quality      string  `json:"quality"`
	LinerScaling float64 `json:"liner_scaling"` // 线性缩放
}

func DriversOperation(listAddress []string, addr string, addressAndValue map[string]DataValue) {
	if len(listAddress) == 6 {
		newValue := addressAndValue[addr].Value
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "int" {
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = float64(addressAndValue[addr].Value.(int)) + value
			}
			if listAddress[4] == "-" {
				newValue = float64(addressAndValue[addr].Value.(int)) - value
			}
			if listAddress[4] == "*" {
				newValue = float64(addressAndValue[addr].Value.(int)) * value
			}
			if listAddress[4] == "/" {
				newValue = float64(addressAndValue[addr].Value.(int)) / value
			}
			if listAddress[4] == "%" {
				newValue = addressAndValue[addr].Value.(int) % int(value)
			}
			if listAddress[4] == "^" {
				newValue = addressAndValue[addr].Value.(int) ^ int(value)
			}
		}
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "uint" {
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = float64(addressAndValue[addr].Value.(uint)) + value
			}
			if listAddress[4] == "-" {
				newValue = float64(addressAndValue[addr].Value.(uint)) - value
			}
			if listAddress[4] == "*" {
				newValue = float64(addressAndValue[addr].Value.(uint)) * value
			}
			if listAddress[4] == "/" {
				newValue = float64(addressAndValue[addr].Value.(uint)) / value
			}
			if listAddress[4] == "%" {
				newValue = addressAndValue[addr].Value.(uint) % uint(value)
			}
			if listAddress[4] == "^" {
				newValue = addressAndValue[addr].Value.(uint) ^ uint(value)
			}
		}
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "int16" {
			// newValue = addressAndValue[addr].Value.(int16)
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = float64(addressAndValue[addr].Value.(int16)) + value
			}
			if listAddress[4] == "-" {
				newValue = float64(addressAndValue[addr].Value.(int16)) - value
			}
			if listAddress[4] == "*" {
				newValue = float64(addressAndValue[addr].Value.(int16)) * value
			}
			if listAddress[4] == "/" {
				newValue = float64(addressAndValue[addr].Value.(int16)) / value
			}
			if listAddress[4] == "%" {
				newValue = addressAndValue[addr].Value.(int16) % int16(value)
			}
			if listAddress[4] == "^" {
				newValue = addressAndValue[addr].Value.(int16) ^ int16(value)
			}
		}
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "int32" {
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = float64(addressAndValue[addr].Value.(int32)) + value
			}
			if listAddress[4] == "-" {
				newValue = float64(addressAndValue[addr].Value.(int32)) - value
			}
			if listAddress[4] == "*" {
				newValue = float64(addressAndValue[addr].Value.(int32)) * value
			}
			if listAddress[4] == "/" {
				newValue = float64(addressAndValue[addr].Value.(int32)) / value
			}
			if listAddress[4] == "%" {
				newValue = addressAndValue[addr].Value.(int32) % int32(value)
			}
			if listAddress[4] == "^" {
				newValue = addressAndValue[addr].Value.(int32) ^ int32(value)
			}
		}
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "int64" {
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = float64(addressAndValue[addr].Value.(int64)) + value
			}
			if listAddress[4] == "-" {
				newValue = float64(addressAndValue[addr].Value.(int64)) - value
			}
			if listAddress[4] == "*" {
				newValue = float64(addressAndValue[addr].Value.(int64)) * value
			}
			if listAddress[4] == "/" {
				newValue = float64(addressAndValue[addr].Value.(int64)) / value
			}
			if listAddress[4] == "%" {
				newValue = addressAndValue[addr].Value.(int64) % int64(value)
			}
			if listAddress[4] == "^" {
				newValue = addressAndValue[addr].Value.(int64) ^ int64(value)
			}
		}
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "float32" {
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = float64(addressAndValue[addr].Value.(float32)) + value
			}
			if listAddress[4] == "-" {
				newValue = float64(addressAndValue[addr].Value.(float32)) - value
			}
			if listAddress[4] == "*" {
				newValue = float64(addressAndValue[addr].Value.(float32)) * value
			}
			if listAddress[4] == "/" {
				newValue = float64(addressAndValue[addr].Value.(float32)) / value
			}
		}
		if reflect.TypeOf(addressAndValue[addr].Value).Name() == "float64" {
			value, _ := strconv.ParseFloat(listAddress[5], 64)
			if listAddress[4] == "+" {
				newValue = addressAndValue[addr].Value.(float64) + value
			}
			if listAddress[4] == "-" {
				newValue = addressAndValue[addr].Value.(float64) - value
			}
			if listAddress[4] == "*" {
				newValue = addressAndValue[addr].Value.(float64) * value
			}
			if listAddress[4] == "/" {
				newValue = addressAndValue[addr].Value.(float64) / value
			}
		}
		addressAndValue[addr] = DataValue{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Value:   newValue,
			Quality: "good",
		}
	}
}
