package main

import (
	"fmt"

	"github.com/feelingsray/Ray-Utils-Go/encode"
)

func main() {
	s := encode.NewSM2Crypt("pyy")
	priv, pub, _ := s.GenerateKeyToMem()
	fmt.Println(priv)
	fmt.Println(pub)
	privateKey, _ := s.ReadPrivateKeyFromMem(priv)
	publicKey, _ := s.ReadPublicKeyFromMem(pub)
	enData, _ := s.Encrypt("你好哈sadfasdfasdfasdfasdfasdfs", publicKey)
	fmt.Println(enData)
	deData, _ := s.Decrypt(enData, privateKey)
	fmt.Println(deData)
}
