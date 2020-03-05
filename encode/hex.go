package encode

import (
	"encoding/hex"
	"fmt"
)

/*
 * 16进制字符串转Byte数组
 */
func HexStringToByte(src string) ([]byte,error){
	var dst []byte
	_, err := fmt.Sscanf(src, "%X", &dst)
	if err != nil {
		return nil,err
	}
	return dst,nil
}

/*
 * Byte数组转16进制字符串
 */
func ByteToHexString(src []byte) string {
	dst := hex.EncodeToString(src)
	return dst
}
