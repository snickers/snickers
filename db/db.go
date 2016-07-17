package db

import "github.com/snickers/snickers/types"
import "github.com/snickers/snickers/db/memory"

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	// Preset methods
	StorePreset(types.Preset) (types.Preset, error)
	RetrievePreset(string) (types.Preset, error)
	UpdatePreset(string, types.Preset) (types.Preset, error)
	GetPresets() ([]types.Preset, error)

	// Job methods
	StoreJob(types.Job) (map[string]types.Job, error)
	RetrieveJob(string) (types.Job, error)
	UpdateJob(string, types.Job) (types.Job, error)
	GetJobs() ([]types.Job, error)

	ClearDatabase() error
}

// GetDatabase returns a handler for the database
func GetDatabase() (DatabaseInterface, error) {
	db, err := memory.GetDatabase()
	return db, err
}
