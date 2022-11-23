package database

import "database/sql"

// GetName is an example that shows you how to query data
func (db *appdbimpl) ValidToken(token string) (bool, error) {
	var userID string
	err := db.c.QueryRow("SELECT UserID From Users WHERE UserID=?", token).Scan(&userID)
	if err == sql.ErrNoRows {
		return false, err
	}
	return true, err
}
