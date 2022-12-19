package database

import (
	"WASA_Photo/service/errorDefinition"
	"WASA_Photo/service/structs"
	"database/sql"
	"errors"
)

// UploadPhoto is a function that uploads a photo to the database
func (db *appdbimpl) UploadPhoto(owner string, filename string) error {
	result, err := db.c.Exec("INSERT INTO Photos (Owner, Filename, Date) VALUES (?, ?, datetime())", owner, filename)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that the user was not banned
		return errorDefinition.ErrNotAdded
	}
	return err
}

// DeletePhoto is a function that deletes a photo from the database
func (db *appdbimpl) DeletePhoto(photoID string) (string, error) {
	var fileIdentifier string
	err := db.c.QueryRow("SELECT Filename FROM Photos WHERE PhotoID = ?", photoID).Scan(&fileIdentifier)
	if err != nil {
		return "", err
	}
	result, err := db.c.Exec("DELETE FROM Photos WHERE PhotoID=?", photoID)
	if err != nil {
		return "", err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return "", err
	} else if affected == 0 {
		// If no rows were affected, it means that the user was not banned
		return "", errorDefinition.ErrPhotoNotFound
	}
	return fileIdentifier, err
}

// GetPhotoOwner is a function that returns the owner of a photo
func (db *appdbimpl) GetPhotoOwner(photoID string) (string, error) {
	var owner string
	err := db.c.QueryRow("SELECT Owner FROM Photos WHERE PhotoID = ? ", photoID).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, err
}

// GetPhoto is a function that returns a photo name from the database
func (db *appdbimpl) GetPhoto(photoID string) (string, error) {
	var filename string
	err := db.c.QueryRow("SELECT Filename FROM Photos WHERE PhotoID = ?", photoID).Scan(&filename)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errorDefinition.ErrPhotoNotFound
	} else if err != nil {
		return "", err
	}
	return filename, err
}

// GetLikes is a function that returns the likes of a photo
func (db *appdbimpl) GetLikes(photoID string, offset int, userRequesting string) ([]structs.Like, error) {
	var likes []structs.Like

	rows, err := db.c.Query("SELECT UserID AS U FROM Likes WHERE PhotoID=? AND (U,?) NOT IN (SELECT * FROM Bans) AND (?,U) NOT IN (SELECT * FROM Bans) LIMIT 30 OFFSET ?", photoID, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var like structs.Like
		err = rows.Scan(&like.UserID)
		if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}
	if len(likes) == 0 {
		return nil, sql.ErrNoRows
	}
	return likes, nil
}

// GetComments is a function that returns the comments of a photo
func (db *appdbimpl) GetComments(photoID string, offset int, userRequesting string) ([]structs.Comment, error) {
	var comments []structs.Comment

	rows, err := db.c.Query("SELECT CommentID AS C, UserName AS U, Text AS T FROM Comments, Users WHERE Users.UserID = Comments.UserID AND PhotoID=? AND (U,?) NOT IN (SELECT * FROM Bans) AND (?,U) NOT IN (SELECT * FROM Bans) ORDER BY CommentID DESC LIMIT 30 OFFSET ?", photoID, userRequesting, userRequesting, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment structs.Comment
		err = rows.Scan(&comment.CommentID, &comment.UserID, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if len(comments) == 0 {
		return nil, sql.ErrNoRows
	}
	return comments, nil
}
