package db

import "github.com/flavioribeiro/snickers/types"

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	StorePreset(types.Preset) map[string]types.Preset
	RetrievePreset(string) types.Preset
	GetPresets() []types.Preset
}
