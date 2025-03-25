package tools

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
)

func AutoShortStr(len int) string {
	chars := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	v4, _ := uuid.NewV4()
	hex := fmt.Sprintf("%x", md5.Sum(v4.Bytes()))
	val, _ := strconv.ParseInt(hex[8:8+8], 16, 0)
	lHexLong := val & 0x3fffffff
	var outChars string
	for j := 0; j < len; j++ {
		outChars += chars[0x0000003D&lHexLong]
		lHexLong >>= 3
	}
	return fmt.Sprintf("%X%s", time.Now().UnixMilli(), outChars)
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

func Cmp(src []string, dest []string) ([]string, []string, []string) {
	msrc := make(map[string]byte)
	mall := make(map[string]byte)
	var set []string
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) {
			l = len(mall)
		} else {
			set = append(set, v)
		}
	}
	for _, v := range set {
		delete(mall, v)
	}
	var added, deleted []string
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted, set
}

func ReverseStrArr(s []string) []string {
	// 对字符串数组进行正向排序
	sort.Strings(s)
	// 反转排序后的字符串数组
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}
