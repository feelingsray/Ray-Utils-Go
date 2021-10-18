package main

import (
    "fmt"
    "github.com/feelingsray/Ray-Utils-Go/rotp"
    "time"
)

func main() {
    mySecretList := []string{"AAAAAAAC","RAYRAY22"}
    // 自定义
    s := rotp.RTOTPCodeWithTime(mySecretList[0],time.Now())
    fmt.Println(s)
    fmt.Println(rotp.RTOTPVerifyWithTime(s,time.Now(),mySecretList))
    // 内部
    s = rotp.RTOTPCodeWithTime(mySecretList[1],time.Now())
    fmt.Println(s)
    fmt.Println(rotp.RTOTPVerifyWithTime(s,time.Now(),nil))

}
