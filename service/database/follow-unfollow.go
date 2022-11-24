package database

// FollowUser is a function that adds userB to the list of users that userA is following
func (db *appdbimpl) FollowUser(userA string, userB string) error {
	_, err := db.c.Exec("INSERT INTO Followers (UserA, UserB) VALUES (?, ?)", userA, userB)
	return err
}

// FollowUser is a function that removes userB from the list of users userA is following
func (db *appdbimpl) UnfollowUser(userA string, userB string) error {
	_, err := db.c.Exec("DELETE FROM Followers WHERE UserA=? AND UserB=?", userA, userB)
	return err
}

// IsFollowing is a function that checks if userA is following userA
func (db *appdbimpl) IsFollowing(userA string, userB string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE UserA=? AND UserB=?", userA, userB)
	err := row.Scan(&count)
	return count > 0, err
}
