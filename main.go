package main

import (
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/server"
)

func main() {
	log := lager.NewLogger("snickers")
	currentDir, _ := os.Getwd()
	config, err := gonfig.FromJsonFile(currentDir + "/config.json")
	if err != nil {
		panic(err)
	}

	// You can register a sink to forward the logs to anywhere.
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	db, err := db.GetDatabase(config)
	if err != nil {
		panic(err)
	}
	port, _ := config.GetString("PORT", "8000")
	snickersServer := server.New(log, config, "tcp", ":"+port, db)
	snickersServer.Start(true)
}
