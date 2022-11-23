package database

// BlockUser is a function that adds userB to the list of blocked users of userA
func (db *appdbimpl) BlockUser(userA string, userB string) error {
	_, err := db.c.Exec("INSERT INTO Bans (UserA, UserB) VALUES (?, ?)", userA, userB)
	return err
}

// UnblockUser is a function that removes userB from the list of blocked users of userA
func (db *appdbimpl) UnblockUser(userA string, userB string) error {
	_, err := db.c.Exec("DELETE FROM Bans WHERE UserA=? AND UserB=?", userA, userB)
	return err
}

// IsBlocked is a function that checks if userB is blocked by userA
func (db *appdbimpl) IsBlocked(userA string, userB string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Bans WHERE UserA=? AND UserB=?", userA, userB)
	err := row.Scan(&count)
	return count > 0, err
}
