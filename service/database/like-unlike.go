package database

import "WASA_Photo/service/errorDefinition"

// Like is a function that adds a like to a photo
func (db *appdbimpl) Like(photoID string, userID string) error {
	_, err := db.c.Exec("INSERT INTO Likes (UserID, PhotoID) VALUES (?,?)", userID, photoID)
	if err != nil {
		return err
	}
	return nil
}

// Unlike is a function that removes a like from a photo
func (db *appdbimpl) Unlike(photoID string, userID string) error {
	result, err := db.c.Exec("DELETE FROM Likes WHERE UserID=? AND PhotoID=?", userID, photoID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return errorDefinition.ErrLikeNotFound
	}
	return nil
}

// HasLiked is a function that checks if a user has liked a photo
func (db *appdbimpl) HasLiked(photoID string, userID string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE UserID=? AND PhotoID=?", userID, photoID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it'll always return a row
	err := row.Scan(&count)
	return count > 0, err
}
