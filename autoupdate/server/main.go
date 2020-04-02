package main

import (
	"github.com/feelingsray/Ray-Utils-Go/tools"
	"net/http"
	"path"
)

func main() {
	releasePath := path.Join(tools.GetAppPath(),"releases")
	releasePath = "/Users/ray/jylink/Ray-Utils-Go/autoupdate/server/releases"
	http.Handle("/", http.FileServer(http.Dir(releasePath)))
	http.ListenAndServe(":9090", nil)
}
