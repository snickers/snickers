package mongo

import (
	"github.com/snickers/snickers/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Database struct that persists configurations
type Database struct {
	session *mgo.Session
	db      *mgo.Database
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
	instance.session = session
	instance.db = session.DB("snickers")
	return instance, nil
}

// ClearDatabase clears the database
func (r *Database) ClearDatabase() error {
	r.db.DropDatabase()
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

// RetrievePreset retrieves one preset from the database
func (r *Database) RetrievePreset(presetName string) (types.Preset, error) {
	c := r.db.C("presets")
	result := types.Preset{}
	err := c.Find(bson.M{"name": presetName}).One(&result)
	return result, err
}

// GetPresets retrieves all presets of the database
func (r *Database) GetPresets() ([]types.Preset, error) {
	results := []types.Preset{}
	c := r.db.C("presets")
	c.Find(nil).All(&results)
	return results, nil
}
