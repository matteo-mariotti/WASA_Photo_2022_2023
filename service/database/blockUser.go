package database

// BlockUser is a function that adds userB to the list of blocked users of userA
func (db *appdbimpl) BlockUser(userA string, userB string) error {
	_, err := db.c.Exec("INSERT INTO Bans (UserA, UserB) VALUES (?, ?)", userA, userB)
	return err
}
