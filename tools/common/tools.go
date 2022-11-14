package common

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/satori/go.uuid"
)

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
