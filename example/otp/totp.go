package main

import (
    "fmt"
    "github.com/feelingsray/Ray-Utils-Go/rotp"
)

func main() {

    s := rotp.RTOTPCode("rayray00")
    fmt.Println(s)
    fmt.Println(rotp.RTOTPVerify(s))

}
