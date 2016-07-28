package main

import (
	"github.com/pivotal-golang/lager"
	"github.com/snickers/snickers/server"
)

func main() {
	log := lager.NewLogger("snickers-server")
	snickersServer := server.New(log, "tcp", ":8080")
	snickersServer.Start()
	select {}
}
