package database

import "database/sql"

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetUsers(start string, offset int) ([]string, error) {
	var userNames []string
	rows, err := db.c.Query("SELECT UserName From Users WHERE UserName LIKE ? LIMIT 30 OFFSET ?", start+"%", offset)
	if err == sql.ErrNoRows {
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
