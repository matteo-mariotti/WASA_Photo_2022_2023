package database

import "WASA_Photo/service/errorDefinition"

// banUser is a function that adds userB to the list of blocked users of userA
func (db *appdbimpl) BanUser(userA string, userB string) error {
	result, err := db.c.Exec("INSERT INTO Bans (UserA, UserB) VALUES (?, ?)", userA, userB)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that I wansn't able to ban the user
		return errorDefinition.ErrUserNotFound
	}
	return nil

}

// UnbanUser is a function that removes userB from the list of blocked users of userA
func (db *appdbimpl) UnbanUser(userA string, userB string) error {
	result, err := db.c.Exec("DELETE FROM Bans WHERE UserA=? AND UserB=?", userA, userB)
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

// IsBanned is a function that checks if userB is blocked by userA
func (db *appdbimpl) IsBanned(userA string, userB string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Bans WHERE UserA=? AND UserB=?", userA, userB)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count > 0, err
}
