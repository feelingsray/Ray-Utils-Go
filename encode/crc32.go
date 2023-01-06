package encode

import "hash/crc32"

func Crc32HashKey(key string, mod uint32) (uint32, uint32) {
	scratch := []byte(key)
	hres := crc32.ChecksumIEEE(scratch)
	mres := hres % mod
	return hres, mres
}
