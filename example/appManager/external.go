package main

import (
    "github.com/feelingsray/Ray-Utils-Go/appManager"
)

func main() {
    ext := appManager.NewExternalManager(9999,3)
    _ = ext.CreateProcStore("mongo","Mongo","/Users/ray/env/mongodb/bin/mongod --config /Users/ray/env/mongodb/mongod.conf",true,true)
    go ext.Manager()

    select {

    }
}
