// Copyright © 2018 Andrea Funtò - released under the MIT License

package main

import (
	"os"

	"github.com/dihedron/go-log"
	"github.com/dihedron/rmq/cmd"
)

func init() {
	log.SetLevel(log.DBG)
	log.SetStream(os.Stdout, true)
	log.SetTimeFormat("15:04:05.000")
	log.SetPrintCallerInfo(true)
	log.SetPrintSourceInfo(log.SourceInfoShort)
}

func main() {
	cmd.Execute()
}
