package db

import (
	"os"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db/memory"
	"github.com/snickers/snickers/db/mongo"
	"github.com/snickers/snickers/types"
)

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	// Preset methods
	StorePreset(types.Preset) (types.Preset, error)
	RetrievePreset(string) (types.Preset, error)
	UpdatePreset(string, types.Preset) (types.Preset, error)
	GetPresets() ([]types.Preset, error)
	DeletePreset(string) (bool, error)

	// Job methods
	StoreJob(types.Job) (types.Job, error)
	RetrieveJob(string) (types.Job, error)
	UpdateJob(string, types.Job) (types.Job, error)
	GetJobs() ([]types.Job, error)

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
