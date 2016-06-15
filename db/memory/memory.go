package memory

// Database struct that persists configurations
type Database struct{}

// NewDatabase creates a new database
func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

//CreatePreset stores preset information in memory
func (r *Database) CreatePreset() int {
	return 0
}
