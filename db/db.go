package db

import "github.com/flavioribeiro/snickers/types"
import "github.com/flavioribeiro/snickers/db/memory"

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	StorePreset(types.Preset) (map[string]types.Preset, error)
	RetrievePreset(string) (types.Preset, error)
	UpdatePreset(string, types.Preset) (types.Preset, error)
	GetPresets() ([]types.Preset, error)
	ClearDatabase() error
}

// GetDatabase returns a handler for the database
func GetDatabase() (DatabaseInterface, error) {
	db, err := memory.GetDatabase()
	return db, err
}
