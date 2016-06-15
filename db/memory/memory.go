package memory

// Database struct that persists configurations
type Database struct{}

// NewDatabase creates a new database
func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

func (r *Database) CreatePreset() int {
	return 0
}
