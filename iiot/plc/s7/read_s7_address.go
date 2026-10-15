package s7

import (
	"log"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type FromMyAddressS7 struct {
	client          *S7Client
	oldAddress      *[]string
	address         map[string][]uint16
	newAddress      map[string][]uint16
	AddressSplit    map[string][][]uint16
	addressAndByte  map[string]map[string][]byte
	AddressAndValue map[string]DataValue
}

func RemoveDuplicates(list []uint16) []uint16 {
	seen := make(map[uint16]bool)
	var uniqueList []uint16
	for _, item := range list {
		if !seen[item] {
			seen[item] = true
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}

func SplitList(list []uint16) [][]uint16 {
	result := make([][]uint16, 0)
	currentSubList := make([]uint16, 0)

	for i := range list {
		currentSubList = append(currentSubList, list[i])

		if i+1 < len(list) && list[i+1]-list[i] > 1 {
			result = append(result, currentSubList)
			currentSubList = make([]uint16, 0)
		}
	}
	if len(currentSubList) > 0 {
		result = append(result, currentSubList)
	}
	return result
}

func ChunkArray(array []uint16, size int) [][]uint16 {
	result := make([][]uint16, 0)
	for i := 0; i < len(array); i += size {
		end := min(i+size, len(array))
		chunk := array[i:end]
		result = append(result, chunk)
	}

	// 移除空的子数组
	for i := 0; i < len(result); i++ {
		if len(result[i]) == 0 {
			result = append(result[:i], result[i+1:]...)
			i--
		}
	}

	return result
}

func round(number float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(number*shift) / shift
}

func InitS7Address(s7Client *S7Client, address *[]string) *FromMyAddressS7 {
	a := map[string][]uint16{"M": {}, "I": {}, "Q": {}, "V": {}}
	n := map[string][]uint16{"M": {}, "I": {}, "Q": {}, "V": {}}
	addressSplit := map[string][][]uint16{"M": {}, "I": {}, "Q": {}, "V": {}}
	for _, add := range *address {
		listAddress := strings.Split(add[1:len(add)-1], ",")
		listAddress[2] = strings.ToUpper(listAddress[2])             // upper
		listAddress[2] = strings.ReplaceAll(listAddress[2], " ", "") // remove spaces
		if strings.ToUpper(listAddress[2][:1]) == "V" {
			if strings.ToUpper(listAddress[2][:2]) == "VB" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["V"] = append(a["V"], uint16(newAdd))
			} else if strings.ToUpper(listAddress[2][:2]) == "VW" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["V"] = append(a["V"], uint16(newAdd))
				a["V"] = append(a["V"], uint16(newAdd)+1)
			} else if strings.ToUpper(listAddress[2][:2]) == "VD" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["V"] = append(a["V"], uint16(newAdd))
				a["V"] = append(a["V"], uint16(newAdd+1))
				a["V"] = append(a["V"], uint16(newAdd+2))
				a["V"] = append(a["V"], uint16(newAdd+3))
			} else {
				countIndex := strings.Split(listAddress[2][1:], ".")
				if len(countIndex) >= 2 {
					newAdd, _ := strconv.Atoi(countIndex[0])
					// index,_ := strconv.Atoi(countIndex[1])
					a["V"] = append(a["V"], uint16(newAdd))
				}
			}
		}
		if strings.ToUpper(listAddress[2][:1]) == "M" {
			if strings.ToUpper(listAddress[2][:2]) == "MB" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["M"] = append(a["M"], uint16(newAdd))
			} else if strings.ToUpper(listAddress[2][:2]) == "MW" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["M"] = append(a["M"], uint16(newAdd))
				a["M"] = append(a["M"], uint16(newAdd+1))
			} else if strings.ToUpper(listAddress[2][:2]) == "MD" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["M"] = append(a["M"], uint16(newAdd))
				a["M"] = append(a["M"], uint16(newAdd+1))
				a["M"] = append(a["M"], uint16(newAdd+2))
				a["M"] = append(a["M"], uint16(newAdd+3))
			} else {
				countIndex := strings.Split(listAddress[2][1:], ".")
				if len(countIndex) >= 2 {
					newAdd, _ := strconv.Atoi(countIndex[0])
					// index,_ := strconv.Atoi(countIndex[1])
					a["M"] = append(a["M"], uint16(newAdd))
				}
			}
		}
		if strings.ToUpper(listAddress[2][:1]) == "I" {
			if strings.ToUpper(listAddress[2][:2]) == "IB" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["I"] = append(a["I"], uint16(newAdd))
			} else if strings.ToUpper(listAddress[2][:2]) == "IW" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["I"] = append(a["I"], uint16(newAdd))
				a["I"] = append(a["I"], uint16(newAdd+1))
			} else if strings.ToUpper(listAddress[2][:2]) == "ID" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["I"] = append(a["I"], uint16(newAdd))
				a["I"] = append(a["I"], uint16(newAdd+1))
				a["I"] = append(a["I"], uint16(newAdd+2))
				a["I"] = append(a["I"], uint16(newAdd+3))
			} else {
				countIndex := strings.Split(listAddress[2][1:], ".")
				if len(countIndex) >= 2 {
					newAdd, _ := strconv.Atoi(countIndex[0])
					// index,_ := strconv.Atoi(countIndex[1])
					a["I"] = append(a["I"], uint16(newAdd))
				}
			}
		}
		if strings.ToUpper(listAddress[2][:1]) == "Q" {
			if strings.ToUpper(listAddress[2][:2]) == "QB" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["Q"] = append(a["Q"], uint16(newAdd))
			} else if strings.ToUpper(listAddress[2][:2]) == "QW" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["Q"] = append(a["Q"], uint16(newAdd))
				a["Q"] = append(a["Q"], uint16(newAdd+1))
			} else if strings.ToUpper(listAddress[2][:2]) == "QD" {
				newAdd, _ := strconv.Atoi(listAddress[2][2:])
				a["Q"] = append(a["Q"], uint16(newAdd))
				a["Q"] = append(a["Q"], uint16(newAdd+1))
				a["Q"] = append(a["Q"], uint16(newAdd+2))
				a["Q"] = append(a["Q"], uint16(newAdd+3))
			} else {
				countIndex := strings.Split(listAddress[2][1:], ".")
				if len(countIndex) >= 2 {
					newAdd, _ := strconv.Atoi(countIndex[0])
					// index,_ := strconv.Atoi(countIndex[1])
					a["Q"] = append(a["Q"], uint16(newAdd))
				}
			}
		}
		if strings.ToUpper(listAddress[2][:2]) == "DB" {
			dbArray := strings.Split(listAddress[2], ".")
			if len(dbArray) >= 2 {
				if _, ok := a[dbArray[0][2:]]; !ok {
					a[dbArray[0][2:]] = make([]uint16, 0)
				}
				if dbArray[1][:3] == "DBX" {
					newAdd, _ := strconv.Atoi(dbArray[1][3:])
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd))
				} else if dbArray[1][:3] == "DBB" {
					newAdd, _ := strconv.Atoi(dbArray[1][3:])
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd))
				} else if dbArray[1][:3] == "DBW" {
					newAdd, _ := strconv.Atoi(dbArray[1][3:])
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd))
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd+1))
				} else if dbArray[1][:3] == "DBD" {
					newAdd, _ := strconv.Atoi(dbArray[1][3:])
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd))
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd+1))
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd+2))
					a[dbArray[0][2:]] = append(a[dbArray[0][2:]], uint16(newAdd+3))
				}
			}
		}
		for key, addressList := range a {
			if len(addressList) >= 1 {
				addressList = RemoveDuplicates(addressList)
				slices.Sort(addressList)
				ones := SplitList(addressList)
				var newOnce [][]uint16
				for _, addrLi := range ones {
					for _, once := range ChunkArray(addrLi, 120) {
						newOnce = append(newOnce, once)
					}
				}
				addressSplit[key] = newOnce
			}
		}
	}
	// log.Println(addressSplit)
	return &FromMyAddressS7{
		client:          s7Client,
		oldAddress:      address,
		address:         a,
		newAddress:      n,
		AddressSplit:    addressSplit,
		addressAndByte:  map[string]map[string][]byte{},
		AddressAndValue: map[string]DataValue{},
	}
}

func (f *FromMyAddressS7) ReadMyAddress() map[string]DataValue {
	if f.address == nil {
		return map[string]DataValue{}
	}
	for key, addList := range f.AddressSplit {
		f.addressAndByte[key] = map[string][]byte{}
		if len(addList) >= 1 {
			dbNum, err := strconv.Atoi(key)
			// log.Println(err, dbNum)
			if err == nil {
				for _, smallAddrList := range addList {
					res, err := f.client.ReadDB(dbNum, int(smallAddrList[0]), len(smallAddrList))
					// log.Println(res)
					if err != nil {
						log.Printf("s7读取报错：%v,%v", res, err)
					}
					for index := range smallAddrList {
						if len(res) >= 1 {
							// newBytes, _ := GetBitBytes(res, index)
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{res[index]}
						} else {
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{}
						}
					}
				}
			} else if key == "V" {
				for _, smallAddrList := range addList {
					res, err := f.client.ReadV(int(smallAddrList[0]), len(smallAddrList)*2)
					if err != nil {
						log.Printf("s7读取报错：%v,%v", res, err)
					}
					for index := range smallAddrList {
						if len(res) >= 1 {
							// newBytes, _ := GetBitBytes(res, index)
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{res[index]}
						} else {
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{}
						}
					}
				}
			} else if key == "M" {
				for _, smallAddrList := range addList {
					res, err := f.client.ReadM(int(smallAddrList[0]), len(smallAddrList))
					if err != nil {
						log.Printf("s7读取报错：%v,%v", res, err)
					}
					for index := range smallAddrList {
						if len(res) >= 1 {
							// newBytes, _ := GetBitBytes(res, index)
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{res[index]}
						} else {
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{}
						}
					}
				}
			} else if key == "I" {
				for _, smallAddrList := range addList {
					res, err := f.client.ReadI(int(smallAddrList[0]), len(smallAddrList))
					// log.Println("I:", smallAddrList, res)
					if err != nil {
						log.Printf("s7读取报错：%v,%v", res, err)
					}
					for index := range smallAddrList {
						if len(res) >= 1 {
							// newBytes, _ := GetBitBytes(res, index)
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{res[index]}
						} else {
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{}
						}
					}
				}
			} else if key == "Q" {
				for _, smallAddrList := range addList {
					res, err := f.client.ReadQ(int(smallAddrList[0]), len(smallAddrList))
					if err != nil {
						log.Printf("s7读取报错：%v,%v", res, err)
					}
					for index := range smallAddrList {
						if len(res) >= 1 {
							// newBytes, _ := GetBitBytes(res, index)
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{res[index]}
						} else {
							f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{}
						}
					}
				}
			} else {
				for _, smallAddrList := range addList {
					for index := range smallAddrList {
						f.addressAndByte[key][strconv.Itoa(int(smallAddrList[index]))] = []byte{}
					}
				}
				// for _, miniAddrList := range addList {
				//	log.Printf("modbus地址填写错误：%v,%v", miniAddrList)
				//	//f.addressAndByte[key] = append(f.addressAndByte[key], []byte{})
				// }
			}
		}
	}
	for _, addr := range *f.oldAddress {
		f.AddressAndValue[addr] = DataValue{}
		listAddress := strings.Split(addr[1:len(addr)-1], ",")
		listAddress[2] = strings.ToUpper(listAddress[2])             // upper
		listAddress[2] = strings.ReplaceAll(listAddress[2], " ", "") // remove spaces
		key1 := ""
		key2 := ""
		count := 0
		var myValue any
		myValue = -9999
		quality := "bad"
		if listAddress[2][:2] == "DB" {
			dbArray := strings.Split(listAddress[2], ".")
			// log.Println("dbArray", dbArray, len(dbArray))
			key1 = dbArray[0][2:]
			if len(dbArray) == 2 {
				// log.Println(dbArray, dbArray[1], dbArray[1][2] == 'W', dbArray[1][1])
				if dbArray[1][2] == 'B' {
					key2 = dbArray[1][3:]
					if listAddress[0][:3] == "int" {
						if len(f.addressAndByte[key1][key2]) >= 1 {
							myValue = ByteToInt16([]byte{0, f.addressAndByte[key1][key2][0]})
							quality = "good"
						}
					}
					if listAddress[0][:4] == "uint" {
						if len(f.addressAndByte[key1][key2]) >= 1 {
							myValue = ByteToUint16([]byte{0, f.addressAndByte[key1][key2][0]})
							quality = "good"
						}
					}
				} else if dbArray[1][2] == 'W' {
					key2 = dbArray[1][3:]
					i, _ := strconv.Atoi(key2)
					str := strconv.Itoa(i + 1)
					if listAddress[0][:3] == "int" {
						if len(f.addressAndByte[key1][key2]) >= 1 {
							myValue = ByteToInt16([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0]})
							// log.Println(myValue)
							quality = "good"
						}
					}
					if listAddress[0][:4] == "uint" {
						if len(f.addressAndByte[key1][key2]) >= 1 {
							myValue = ByteToUint16(f.addressAndByte[key1][key2])
							quality = "good"
						}
					}
				} else if dbArray[1][2] == 'D' {
					key2 = dbArray[1][3:]
					// log.Println(key2)
					i, _ := strconv.Atoi(key2)
					str := strconv.Itoa(i + 1)
					str2 := strconv.Itoa(i + 2)
					str3 := strconv.Itoa(i + 3)
					if listAddress[0][:3] == "int" {
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 && len(f.addressAndByte[key1][str2]) >= 1 && len(f.addressAndByte[key1][str3]) >= 1 {
							myValue = ByteToInt32([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0], f.addressAndByte[key1][str2][0], f.addressAndByte[key1][str3][0]}, false)
							quality = "good"
						}
					}
					if listAddress[0][:4] == "uint" {
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 && len(f.addressAndByte[key1][str2]) >= 1 && len(f.addressAndByte[key1][str3]) >= 1 {
							myValue = ByteToUint32([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0], f.addressAndByte[key1][str2][0], f.addressAndByte[key1][str3][0]}, false)
							quality = "good"
						}
					}
					if listAddress[0][:5] == "float" {
						// str := "123"
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 && len(f.addressAndByte[key1][str2]) >= 1 && len(f.addressAndByte[key1][str3]) >= 1 {
							myValue = ByteToFloat32([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0], f.addressAndByte[key1][str2][0], f.addressAndByte[key1][str3][0]}, false)
							quality = "good"
						}
					}
				}
			} else if len(dbArray) == 3 {
				key2 = dbArray[1][3:]
				count, _ = strconv.Atoi(dbArray[2])
				if dbArray[1][2] == 'X' {
					var newByte byte
					if len(f.addressAndByte[key1][key2]) >= 1 {
						newByte = f.addressAndByte[key1][key2][0]
					}
					if len(f.addressAndByte[key1][key2]) >= 1 {
						myValue = GetBit(newByte, uint(count))
						quality = "good"
					}
				}
			}
		} else {
			dbArray := strings.Split(listAddress[2], ".")
			key1 = dbArray[0][:1]
			if len(dbArray) == 1 {
				key2 = dbArray[0][2:]
				if dbArray[0][1] == 'B' {
					if listAddress[0][:3] == "int" {
						if len(f.addressAndByte[key1][key2]) >= 1 {
							myValue = ByteToInt16([]byte{0, f.addressAndByte[key1][key2][0]})
							quality = "good"
						}
					}
					if listAddress[0][:4] == "uint" {
						if len(f.addressAndByte[key1][key2]) >= 1 {
							myValue = ByteToUint16(f.addressAndByte[key1][key2])
							quality = "good"
						}
					}
					// value = ByteToUint16(f.addressAndByte[key1][key2])
				} else if listAddress[2][1] == 'W' {
					i, _ := strconv.Atoi(key2)
					str := strconv.Itoa(i + 1)

					if listAddress[0][:3] == "int" {
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 {
							myValue = ByteToInt16([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0]})
							quality = "good"
						}
					}
					if listAddress[0][:4] == "uint" {
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 {
							myValue = ByteToUint16([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0]})
							quality = "good"
						}
					}
				} else if listAddress[2][1] == 'D' {
					i, _ := strconv.Atoi(key2)
					str := strconv.Itoa(i + 1)
					str2 := strconv.Itoa(i + 2)
					str3 := strconv.Itoa(i + 3)
					if listAddress[0][:3] == "int" {
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 && len(f.addressAndByte[key1][str2]) >= 1 && len(f.addressAndByte[key1][str3]) >= 1 {
							myValue = ByteToInt32([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0], f.addressAndByte[key1][str2][0], f.addressAndByte[key1][str3][0]}, false)
							quality = "good"
						}
					}
					if listAddress[0][:4] == "uint" {
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 && len(f.addressAndByte[key1][str2]) >= 1 && len(f.addressAndByte[key1][str3]) >= 1 {
							myValue = ByteToUint32([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0], f.addressAndByte[key1][str2][0], f.addressAndByte[key1][str3][0]}, false)
							quality = "good"
						}
					}
					if listAddress[0][:5] == "float" {
						// str := "123"
						if len(f.addressAndByte[key1][key2]) >= 1 && len(f.addressAndByte[key1][str]) >= 1 && len(f.addressAndByte[key1][str2]) >= 1 && len(f.addressAndByte[key1][str3]) >= 1 {
							myValue = ByteToFloat32([]byte{f.addressAndByte[key1][key2][0], f.addressAndByte[key1][str][0], f.addressAndByte[key1][str2][0], f.addressAndByte[key1][str3][0]}, false)
							quality = "good"
						}
					}
				}
			} else if len(dbArray) == 2 {
				count, _ = strconv.Atoi(dbArray[1])
				key2 = dbArray[0][1:]
				var newByte byte
				if len(f.addressAndByte[key1][key2]) == 1 {
					newByte = f.addressAndByte[key1][key2][0]
				}
				if len(f.addressAndByte[key1][key2]) >= 1 {
					myValue = GetBit(newByte, uint(count))
					quality = "good"
				}
			}
		}
		f.AddressAndValue[addr] = DataValue{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Value:   myValue,
			Quality: quality,
		}

		if myValue != -9999 {
			if len(listAddress) == 6 {
				news := f.AddressAndValue
				DriversOperation(listAddress, addr, news)
			}
		}
	}
	return f.AddressAndValue
}
