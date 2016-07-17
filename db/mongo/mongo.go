package mongo

import (
	"github.com/snickers/snickers/types"
	"gopkg.in/mgo.v2"
)

// Database struct that persists configurations
type Database struct {
	db *mgo.Database
}

var instance *Database

// GetDatabase returns database singleton
func GetDatabase() (*Database, error) {
	instance = &Database{}
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	instance.db = session.DB("snickers")
	return instance, nil
}

// ClearDatabase clears the database
func (r *Database) ClearDatabase() error {
	return nil
}

// StorePreset stores preset information
func (r *Database) StorePreset(preset types.Preset) (types.Preset, error) {
	c := r.db.C("presets")
	err := c.Insert(preset)
	if err != nil {
		return types.Preset{}, err
	}
	return preset, nil
}
