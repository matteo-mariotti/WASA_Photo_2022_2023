package database

import "database/sql"

// ValidToken checks if a token is valid
func (db *appdbimpl) ValidToken(token string) (bool, error) {
	var userID string
	err := db.c.QueryRow("SELECT UserID From Users WHERE UserID=?", token).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return true, err
}
