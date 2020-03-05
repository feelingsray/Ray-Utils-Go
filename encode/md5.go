package encode

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Str(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
