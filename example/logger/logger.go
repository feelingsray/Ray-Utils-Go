package main

import "github.com/feelingsray/Ray-Utils-Go/logger"

func main() {

	log:= logger.LoggerConsoleHandle(logger.DebugLevel)
	log.Info("OK");
}
