package s7

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
)

func Float32ToByte(float float32, swap bool) []byte {
	if swap {
		bits := math.Float32bits(float)
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits>>16|bits<<16)
		return bytes
	}
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)
	return bytes
}

func ByteToFloat32(bytes []byte, swap bool) float32 {
	if swap {
		// newBytes := []byte{bytes[2], bytes[3], bytes[0], bytes[1]}
		bits := binary.BigEndian.Uint32(bytes)
		// log.Println(bits)
		return math.Float32frombits(bits>>16 | bits<<16)
	}
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToByte(float float64, swap bool) []byte {
	if swap {
		bits := math.Float64bits(float)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits>>32|bits<<32)
		return bytes
	}
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

func ByteToFloat64(bytes []byte, swap bool) float64 {
	if swap {
		// newBytes := []byte{bytes[6], bytes[7], bytes[4], bytes[5], bytes[2], bytes[3], bytes[0], bytes[1]}
		bits := binary.BigEndian.Uint64(bytes)
		return math.Float64frombits(bits>>32 | bits<<32)
	}
	bits := binary.BigEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func ByteToInt32(bytes []byte, swap bool) int32 {
	if swap {
		// newBytes := []byte{bytes[6], bytes[7], bytes[4], bytes[5], bytes[2], bytes[3], bytes[0], bytes[1]}
		bits := binary.BigEndian.Uint32(bytes)
		return int32(bits>>16 | bits<<16)
	}
	bits := binary.BigEndian.Uint32(bytes)
	return int32(bits)
}

func ByteToUint32(bytes []byte, swap bool) uint32 {
	if swap {
		bits := binary.BigEndian.Uint32(bytes)
		return bits>>16 | bits<<16
	}
	bits := binary.BigEndian.Uint32(bytes)
	return bits
}

func ByteToUint64(bytes []byte, swap bool) uint64 {
	if swap {
		bits := binary.BigEndian.Uint64(bytes)
		return bits>>32 | bits<<32
	}
	bits := binary.BigEndian.Uint64(bytes)
	return bits
}

func Uint32ToByte(uint uint32, swap bool) []byte {
	if swap {
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint>>16|uint<<16)
		return bytes
	}
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint)
	return bytes
}

func Int32ToByte(i int32, swap bool) []byte {
	newI := uint32(i)
	return Uint32ToByte(newI, swap)
}

func Uint64ToByte(uint uint64, swap bool) []byte {
	if swap {
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, uint>>32|uint<<32)
		return bytes
	}
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint)
	return bytes
}

func Int64ToByte(i int64, swap bool) []byte {
	newI := uint64(i)
	return Uint64ToByte(newI, swap)
}

func ByteToUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func Uint16ToByte(uint uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint)
	return bytes
}

func ByteToInt16(bytes []byte) int16 {
	u := ByteToUint16(bytes)
	return int16(u)
}

func Int16ToByte(int int16) []byte {
	u := uint16(int)
	return Uint16ToByte(u)
}

func Int8ToByte(int int8) []byte {
	u := uint16(int)
	return []byte{Uint16ToByte(u)[1]}
}

func GetBit(byteVal byte, n uint) uint {
	// 将 byte 值向右移动 n 位，并与 1 进行按位与操作
	// 结果非零表示第 n 位为 1，否则为 0
	if (byteVal>>n)&1 != 0 {
		return 1
	}
	return 0
}

// GetBit 函数获取 []byte 中指定位的 bit 值
func GetBitBytes(byteSlice []byte, bitPosition int) ([]byte, error) {
	if bitPosition < 0 || bitPosition >= 8*len(byteSlice) {
		return []byte{}, errors.New("Bit position out of range")
	}
	// 获取指定索引的字节
	// byteIndex := bitPosition / 8
	// log.Println(bitPosition)
	// byteIndex := len(byteSlice) - 1 - bitPosition/8 // 字节索引位置
	// bitIndex := uint(7 - (bitPosition % 8))         // 位索引位置
	byteIndex := bitPosition / 8      // 字节索引位置
	bitIndex := uint(bitPosition % 8) // 位索引位置
	// log.Println(byteIndex, bitIndex)
	// mask := byte(1 << uint(7-bitIndex)) // 位掩码
	//
	// value := byteSlice[byteIndex] & mask
	// bitValue := value >> uint(7-bitIndex)
	// 提取比特位的值
	bitValue := (byteSlice[byteIndex] >> bitIndex) & 1
	return []byte{bitValue}, nil
}

func GetBitByInt16(int int16) int {
	return 0
}

func EditBitInt16(num int16, n, value int16) int16 {
	// 将 int16 转换为二进制字符串
	binaryStr := fmt.Sprintf("%016s", strconv.FormatInt(int64(num), 2))
	// 找到第 n 位的位置
	bitIndex := 15 - n // 从 0 开始计数
	if bitIndex >= int16(len(binaryStr)) {
		// 如果 n 超出二进制位数，不做任何修改
		return num
	}
	// 修改第 n 位的值
	newBinaryStr := replaceCharAt(binaryStr, bitIndex, strconv.Itoa(int(value)))
	// 将修改后的二进制字符串转换回 int16
	result, _ := strconv.ParseInt(newBinaryStr, 2, 16)
	return int16(result)
}

func replaceCharAt(s string, index int16, r string) string {
	if index < 0 || index >= int16(len(s)) {
		return s
	}
	return s[:index] + r + s[index+1:]
}
