package database

import "WASA_Photo/service/errorDefinition"

// TODO Comment this function
func (db *appdbimpl) Like(photoID string, userID string) error {
	_, err := db.c.Exec("INSERT INTO Likes (UserID, PhotoID) VALUES (?,?)", userID, photoID)
	if err != nil {
		return err
	}
	return nil
}

// TODO Comment this function
func (db *appdbimpl) Unlike(photoID string, userID string) error {
	result, err := db.c.Exec("DELETE FROM Likes WHERE UserID=? AND PhotoID=?", userID, photoID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that the request is wrong
		return errorDefinition.ErrLikeNotFound
	}
	return nil
}

func (db *appdbimpl) HasLiked(photoID string, userID string) (bool, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE UserID=? AND PhotoID=?", userID, photoID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it'll always return a row
	err := row.Scan(&count)
	return count > 0, err
}
