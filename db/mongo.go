package db

import (
	"sync"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Database struct that persists configurations
type mongoDatabase struct {
	db *mgo.Database
}

var databaseMongoInit sync.Once
var mongoInstance *mongoDatabase

// GetDatabase returns database singleton
func getMongoDatabase(cfg gonfig.Gonfig) (Storage, error) {
	var initErr error
	databaseMongoInit.Do(func() {
		mongoInstance = &mongoDatabase{}
		mongoHost, _ := cfg.GetString("MONGODB_HOST", "")
		session, err := mgo.Dial(mongoHost)
		if err != nil {
			initErr = err
		}
		session.SetMode(mgo.Monotonic, true)
		mongoInstance.db = session.DB("snickers")
	})

	return mongoInstance, initErr
}

// ClearDatabase clears the database
func (r *mongoDatabase) ClearDatabase() error {
	return r.db.DropDatabase()
}

// StorePreset stores preset information
func (r *mongoDatabase) StorePreset(preset types.Preset) (types.Preset, error) {
	c := r.db.C("presets")
	err := c.Insert(preset)
	if err != nil {
		return types.Preset{}, err
	}
	return preset, nil
}

// UpdatePreset updates a preset
func (r *mongoDatabase) UpdatePreset(presetName string, newPreset types.Preset) (types.Preset, error) {
	c := r.db.C("presets")
	err := c.Update(bson.M{"name": presetName}, newPreset)
	if err != nil {
		return types.Preset{}, err
	}
	return newPreset, nil
}

// RetrievePreset retrieves one preset from the database
func (r *mongoDatabase) RetrievePreset(presetName string) (types.Preset, error) {
	c := r.db.C("presets")
	result := types.Preset{}
	err := c.Find(bson.M{"name": presetName}).One(&result)
	return result, err
}

// GetPresets retrieves all presets of the database
func (r *mongoDatabase) GetPresets() ([]types.Preset, error) {
	results := []types.Preset{}
	c := r.db.C("presets")
	err := c.Find(nil).All(&results)
	return results, err
}

// DeletePreset deletes a preset from the database
func (r *mongoDatabase) DeletePreset(presetName string) (types.Preset, error) {
	result, err := r.RetrievePreset(presetName)
	if err != nil {
		return types.Preset{}, err
	}

	c := r.db.C("presets")
	err = c.Remove(bson.M{"name": presetName})
	if err != nil {
		return types.Preset{}, err
	}
	return result, nil
}

// StoreJob stores job information
func (r *mongoDatabase) StoreJob(job types.Job) (types.Job, error) {
	c := r.db.C("jobs")
	err := c.Insert(job)
	if err != nil {
		return types.Job{}, err
	}
	return job, nil
}

// RetrieveJob retrieves one job from the database
func (r *mongoDatabase) RetrieveJob(jobID string) (types.Job, error) {
	c := r.db.C("jobs")
	result := types.Job{}
	err := c.Find(bson.M{"id": jobID}).One(&result)
	return result, err
}

// UpdateJob updates a job
func (r *mongoDatabase) UpdateJob(jobID string, newJob types.Job) (types.Job, error) {
	c := r.db.C("jobs")
	err := c.Update(bson.M{"id": jobID}, newJob)
	if err != nil {
		return types.Job{}, err
	}
	return newJob, nil
}

//GetJobs retrieves all jobs of the database
func (r *mongoDatabase) GetJobs() ([]types.Job, error) {
	results := []types.Job{}
	c := r.db.C("jobs")
	err := c.Find(nil).All(&results)
	return results, err
}
