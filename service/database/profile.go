package database

func (db *appdbimpl) GetFollowerNumber(userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE UserB=?", userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}
func (db *appdbimpl) GetFollowingNumber(userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE UserA=?", userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}
func (db *appdbimpl) GetPhotosNumber(userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Photos WHERE Owner=?", userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}

func (db *appdbimpl) GetPhotos(userID string, offset int) ([]int, error) {
	var photos []int
	rows, err := db.c.Query("SELECT PhotoID FROM Photos WHERE UserID=? LIMIT 30 OFFSET ?", userID, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var photoID int
		err = rows.Scan(&photoID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photoID)
	}
	return photos, nil
}
