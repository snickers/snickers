package db

import "github.com/flavioribeiro/snickers/types"
import "github.com/flavioribeiro/snickers/db/memory"

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	StorePreset(types.Preset) map[string]types.Preset
	RetrievePreset(string) types.Preset
	GetPresets() []types.Preset
	ClearDatabase()
}

// GetDatabase returns a handler for the database
func GetDatabase() (DatabaseInterface, error) {
	db, err := memory.GetDatabase()
	return db, err
}
