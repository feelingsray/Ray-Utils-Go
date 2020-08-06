package main

import (
	"fmt"
	"github.com/feelingsray/Ray-Utils-Go/encode"
)

func main() {

	s := encode.NewSM4Crypt("1234567890123457")
	enData, err := s.Encrypt("asdfasdfasdfasdfasdfasdfasdfasdf")
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println(enData)
	deData, err := s.Decrypt(enData)
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println(deData)
}
