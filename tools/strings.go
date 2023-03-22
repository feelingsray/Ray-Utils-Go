package tools

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	
	uuid "github.com/satori/go.uuid"
)

func AutoShortStr(len int) string {
	chars := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	v4, _ := uuid.NewV4()
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

func StringsContains(obj string, list []string) bool {
	for _, v := range list {
		if v == obj {
			return true
		}
	}
	return false
}

func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	lenStr := len(runeStr)
	
	if start < 0 {
		start = lenStr + start
	}
	if start > lenStr {
		start = lenStr
	}
	end := start + length
	if end > lenStr {
		end = lenStr
	}
	if length < 0 {
		end = lenStr + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}
