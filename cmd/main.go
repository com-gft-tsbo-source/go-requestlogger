package main

import (
	"os"

	"github.com/com-gft-tsbo-source/go-requestlogger/requestlogger"
)

// ###########################################################################
// ###########################################################################
// MAIN
// ###########################################################################
// ###########################################################################

var usage []byte = []byte("requestlogger: [OPTIONS] ")

func main() {

	var ms requestlogger.RequestLogger
	requestlogger.InitFromArgs(&ms, os.Args, nil)
	ms.Run()
}
