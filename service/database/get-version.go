package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetVersion() (string, error) {
	var version string
	err := db.c.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	return version, err
}
