package database

import "WASA_Photo/service/errorDefinition"

// FollowUser is a function that adds userB to the list of users that userA is following
func (db *appdbimpl) FollowUser(userA string, userB string) error {
	result, err := db.c.Exec("INSERT INTO Followers (UserA, UserB) VALUES (?, ?)", userA, userB)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that the user was not banned
		return errorDefinition.ErrUserNotFound
	}
	return nil
}

// FollowUser is a function that removes userB from the list of users userA is following
func (db *appdbimpl) UnfollowUser(userA string, userB string) error {
	result, err := db.c.Exec("DELETE FROM Followers WHERE UserA=? AND UserB=?", userA, userB)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that the user was not banned
		return errorDefinition.ErrUserNotFound
	}
	return nil
}

// IsFollowing is a function that checks if userA is following userA
func (db *appdbimpl) IsFollowing(userA string, userB string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE UserA=? AND UserB=?", userA, userB)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count > 0, err
}

// GetFollowing is a function that returns a list of users that userID is following, internal use only, no pagination
func (db *appdbimpl) GetFollowing(userID string) ([]string, error) {
	var following []string
	rows, err := db.c.Query("SELECT UserB FROM Followers WHERE UserA=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follow string
		err := rows.Scan(&follow)
		if err != nil {
			return nil, err
		}
		following = append(following, follow)
	}
	return following, nil
}
