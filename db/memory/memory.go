package memory

import "github.com/flavioribeiro/snickers/db"

// Database struct that persists configurations
type Database struct {
	Presets map[string]db.Preset
}

// NewDatabase creates a new database
func NewDatabase() (*Database, error) {
	d := &Database{}
	d.Presets = map[string]db.Preset{}
	return d, nil
}

//CreatePreset stores preset information
func (r *Database) StorePreset(preset db.Preset) map[string]db.Preset {
	r.Presets[preset.Name] = preset
	return r.Presets
}

//RetrievePreset retrieves one preset from the database
func (r *Database) RetrievePreset(presetName string) db.Preset {
	return r.Presets[presetName]
}
