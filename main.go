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
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

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

	port, err := config.GetString("PORT", "8000")
	if err != nil {
		panic(err)
	}

	snickersServer := server.New(log, config, "tcp", ":"+port, db)
	snickersServer.Start(true)
}
