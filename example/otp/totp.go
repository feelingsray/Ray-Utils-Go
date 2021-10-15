package main

import (
    "fmt"
    "github.com/feelingsray/Ray-Utils-Go/rotp"
    "time"
)

func main() {

    d,_ := time.ParseDuration("-1h")
    e,_ := time.ParseDuration("-2h")
    _ = e
    s := rotp.RTOTPCodeWithTime("rayray00",time.Now().Add(d))
    fmt.Println(s)
    fmt.Println(rotp.RTOTPVerifyWithTime(s,time.Now().Add(d)))

}
