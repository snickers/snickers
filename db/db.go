package db

import (
	"os"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db/memory"
	"github.com/snickers/snickers/db/mongo"
	"github.com/snickers/snickers"
)

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	// Preset methods
	StorePreset(snickers.Preset) (snickers.Preset, error)
	RetrievePreset(string) (snickers.Preset, error)
	UpdatePreset(string, snickers.Preset) (snickers.Preset, error)
	GetPresets() ([]snickers.Preset, error)
	DeletePreset(string) (snickers.Preset, error)

	// Job methods
	StoreJob(snickers.Job) (snickers.Job, error)
	RetrieveJob(string) (snickers.Job, error)
	UpdateJob(string, snickers.Job) (snickers.Job, error)
	GetJobs() ([]snickers.Job, error)

	ClearDatabase() error
}

// GetDatabase returns a handler for the database
func GetDatabase() (DatabaseInterface, error) {
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	mongoHost, _ := cfg.GetString("MONGODB_HOST", "")

	if mongoHost != "" {
		return mongo.GetDatabase()
	}
	return memory.GetDatabase()
}
