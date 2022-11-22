package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) LoginUser(username string) (string, error) {
	var uuid string

	row := db.c.QueryRow("SELECT UserID FROM Users WHERE UserName=?", username)
	err := row.Scan(&uuid)
	return uuid, err
}
