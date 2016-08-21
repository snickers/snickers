package db

import "github.com/snickers/snickers/types"

//go:generate counterfeiter . Storage

// Storage defines functions for accessing data
type Storage interface {
	// Preset methods
	StorePreset(types.Preset) (types.Preset, error)
	RetrievePreset(string) (types.Preset, error)
	UpdatePreset(string, types.Preset) (types.Preset, error)
	GetPresets() ([]types.Preset, error)
	DeletePreset(string) (types.Preset, error)

	// Job methods
	StoreJob(types.Job) (types.Job, error)
	RetrieveJob(string) (types.Job, error)
	UpdateJob(string, types.Job) (types.Job, error)
	GetJobs() ([]types.Job, error)

	ClearDatabase() error
}
