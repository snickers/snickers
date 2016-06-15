package memory

import "github.com/flavioribeiro/snickers/db"

// Database struct that persists configurations
type Database struct {
	Presets []db.Preset
}

// NewDatabase creates a new database
func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

//CreatePreset stores preset information in memory
func (r *Database) CreatePreset(preset db.Preset) []db.Preset {
	r.Presets = append(r.Presets, preset)
	return r.Presets
}
