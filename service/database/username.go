package database

// UserAvailable checks if a username is available
func (db *appdbimpl) UserAvailable(username string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT Count(*) FROM Users WHERE UserName=?", username)
	err := row.Scan(&count)
	return count > 0, err
}

// ChangeUsername changes the username of a user
func (db *appdbimpl) ChangeUsername(userID string, newname string) error {
	_, err := db.c.Exec("UPDATE Users SET UserName=? WHERE UserID=?", newname, userID)
	return err
}

// GetName gets the username of a user
func (db *appdbimpl) GetName(userID string) (string, error) {
	var username string
	row := db.c.QueryRow("SELECT UserName FROM Users WHERE UserID=?", userID)
	err := row.Scan(&username)
	return username, err
}
