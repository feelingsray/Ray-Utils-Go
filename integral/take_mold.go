package integral

import "crypto/md5"

func TakeMold(raw string, num int) (mold int) {
	data := []byte(raw)
	md5Str := md5.Sum(data)
	mold = int(md5Str[0]) % num
	return mold
}
