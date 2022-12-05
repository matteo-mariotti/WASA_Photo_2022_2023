package database

// LoginUser get the userID from the database
func (db *appdbimpl) LoginUser(username string) (string, error) {
	var uuid string

	row := db.c.QueryRow("SELECT UserID FROM Users WHERE UserName=?", username)
	err := row.Scan(&uuid)
	return uuid, err
}

// RegisterUser adds a new user to the database
func (db *appdbimpl) RegisterUser(uuid string, username string) error {
	_, err := db.c.Exec("INSERT INTO Users (UserID, UserName) VALUES (?, ?)", uuid, username)
	return err
}
