package encode

// BytesToBinaryString get the string in binary format of a []byte or []int8.
func BytesToBinaryString(bs []byte) string {
	l := len(bs)
	bl := l*8
	buf := make([]byte, 0, bl)
	for _, b := range bs {
		buf = appendBinaryString(buf, b)
	}
	return string(buf)
}

// append bytes of string in binary format.
func appendBinaryString(bs []byte, b byte) []byte {
	var a byte
	for i := 0; i < 8; i++ {
		a = b
		b <<= 1
		b >>= 1
		switch a {
		case b:
			bs = append(bs, byte('0'))
		default:
			bs = append(bs, byte('1'))
		}
		b <<= 1
	}
	return bs
}
