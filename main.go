package main

import (
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/snickers/snickers/server"
)

func main() {
	log := lager.NewLogger("snickers")
	// You can register a sink to foward the logs to anywhere.
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	snickersServer := server.New(log, "tcp", ":8080")
	snickersServer.Start(true)
}
