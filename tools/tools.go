package tools

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/satori/go.uuid"
)

func ReverseInterfaceArray(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func ReverseInts(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func ReverseStrings(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func ReverseFloats(s []float64) []float64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func ReverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	rune_str := []rune(str)
	len_str := len(rune_str)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(rune_str[start:end])
}

// 短字符串
func HashShortStr(data string) []string {
	chars := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	hex := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	resUrl := make([]string, 4)
	for i := 0; i < 4; i++ {
		val, _ := strconv.ParseInt(hex[i*8:i*8+8], 16, 0)
		lHexLong := val & 0x3fffffff
		outChars := ""
		for j := 0; j < 6; j++ {
			outChars += chars[0x0000003D&lHexLong]
			lHexLong >>= 5
		}
		resUrl[i] = outChars
	}
	return resUrl
}

func AutoShortStr(len int) string {
	chars := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	v4 := uuid.NewV4()
	hex := fmt.Sprintf("%x", md5.Sum(v4.Bytes()))
	val, _ := strconv.ParseInt(hex[8:8+8], 16, 0)
	lHexLong := val & 0x3fffffff
	outChars := ""
	for j := 0; j < len; j++ {
		outChars += chars[0x0000003D&lHexLong]
		lHexLong >>= 3
	}
	return outChars
}

func ObjectInArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func StringsContains(obj string, list []string) bool {
	for _, v := range list {
		if v == obj {
			return true
		}
	}
	return false
}
