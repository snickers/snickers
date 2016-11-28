package main

import (
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/server"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config, err := gonfig.FromJsonFile(currentDir + "/config.json")
	if err != nil {
		panic(err)
	}

	log, err := setupLogger(config)
	if err != nil {
		panic(err)
	}

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

func setupLogger(config gonfig.Gonfig) (lager.Logger, error) {
	log := lager.NewLogger("snickers")
	logfile, err := config.GetString("LOGFILE", "")
	if err != nil {
		return nil, err
	}
	if logfile != "" {
		f, err := os.Create(logfile)
		if err != nil {
			return nil, err
		}
		log.RegisterSink(lager.NewWriterSink(f, lager.DEBUG))
	} else {
		log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	}
	return log, nil
}
