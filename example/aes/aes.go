package main

import (
	"fmt"

	"github.com/feelingsray/Ray-Utils-Go/encode"
)

func main() {

	KEY := "7Z/Wj7LhXemfsfIHaKcTtA=="
	IV := "ZMta/6i+v2z/NZmQjIWAiA=="

	obj := encode.NewAESCryptor(KEY, IV)
	enData, _ := obj.AESBase64Encrypt("ssssss")
	fmt.Println(enData)
	data, _ := obj.AESBase64Decrypt(enData)
	fmt.Println(data)

}
