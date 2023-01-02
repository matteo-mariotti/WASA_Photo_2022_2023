package database

import (
	"database/sql"
	"errors"
)

// GetUsers is a function that returns a list of users matching the given search string
func (db *appdbimpl) GetUsers(start string, offset int, reqUser string) ([]string, error) {
	var userNames []string
	rows, err := db.c.Query("SELECT UserName From Users WHERE UserName LIKE ? AND (Users.UserID, ?) NOT IN (SELECT * FROM Bans) LIMIT 30 OFFSET ?", start+"%", reqUser, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	for rows.Next() {
		var userName string
		err = rows.Scan(&userName)
		if err != nil {
			return nil, err
		}
		userNames = append(userNames, userName)
	}
	if len(userNames) == 0 {
		return nil, sql.ErrNoRows
	}

	return userNames, err
}
