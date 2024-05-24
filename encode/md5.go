package encode

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5Bytes(s []byte) string {
	ret := md5.Sum(s)
	return hex.EncodeToString(ret[:])
}

// MD5 计算字符串MD5值
func MD5(s string) string {
	return MD5Bytes([]byte(s))
}

// MD5File 计算文件MD5值
func MD5File(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	return MD5Bytes(data)
}
