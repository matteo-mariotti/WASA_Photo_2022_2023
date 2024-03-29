package database

import "WASA_Photo/service/errorDefinition"

// Comment is a function that adds a comment to a photo
func (db *appdbimpl) Comment(photoID string, userID string, text string) error {
	_, err := db.c.Exec("INSERT INTO Comments (PhotoID, UserID, Text) VALUES (?, ?, ?)", photoID, userID, text)
	if err != nil {
		return err
	}
	return nil
}

// Uncomment is a function that removes a comment from a photo
func (db *appdbimpl) Uncomment(photoID string, userID string, commentID string) error {
	result, err := db.c.Exec("DELETE FROM Comments WHERE CommentID=? AND UserID=? AND PhotoID=? ", commentID, userID, photoID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that the request is wrong
		return errorDefinition.ErrCommmentNotFound
	}
	return nil
}
