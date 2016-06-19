package db

// DatabaseInterface defines functions for accessing data
type DatabaseInterface interface {
	StorePreset(Preset) map[string]Preset
	RetrievePreset(string) Preset
	GetPresets() []Preset
}
