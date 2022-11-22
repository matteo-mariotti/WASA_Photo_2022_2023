package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) RegisterUser(uuid string, username string) error {
	_, err := db.c.Exec("INSERT INTO Users (UserID, UserName) VALUES (?, ?)", uuid, username)
	return err
}
