package main

import (
    "fmt"
    "github.com/feelingsray/Ray-Utils-Go/rotp"
    "time"
)

func main() {
    s := rotp.RTOTPCodeWithTime("RAYRAY22",time.Now())
    fmt.Println(s)
    fmt.Println(rotp.RTOTPVerifyWithTime(s,time.Now()))

}
