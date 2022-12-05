package database

import (
	"WASA_Photo/service/errorDefinition"
	"WASA_Photo/service/structs"
	"database/sql"
)

// TODO Comment
func (db *appdbimpl) UploadPhoto(owner string, filename string) error {
	result, err := db.c.Exec("INSERT INTO Photos (Owner, Filename, Date) VALUES (?, ?, datetime())", owner, filename)

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If no rows were affected, it means that the user was not banned
		return errorDefinition.ErrNotAdded
	}
	return err
}

// TODO Comment
func (db *appdbimpl) DeletePhoto(photoID string) (string, error) {
	var fileIdentifier string
	err := db.c.QueryRow("SELECT Filename FROM Photos WHERE PhotoID = ?", photoID).Scan(&fileIdentifier)
	if err != nil {
		return "", err
	}
	result, err := db.c.Exec("DELETE FROM Photos WHERE PhotoID=?", photoID)

	affected, err := result.RowsAffected()
	if err != nil {
		return "", err
	} else if affected == 0 {
		// If no rows were affected, it means that the user was not banned
		return "", errorDefinition.ErrPhotoNotFound
	}
	return fileIdentifier, err
}

// TODO Comment
func (db *appdbimpl) GetPhotoOwner(photoID string) (string, error) {
	var owner string
	err := db.c.QueryRow("SELECT Owner FROM Photos WHERE PhotoID = ?", photoID).Scan(&owner)
	if err != nil {
		return "", err
	}
	return owner, err
}

// TODO Comment
func (db *appdbimpl) GetPhoto(photoID string) (string, error) {
	var filename string
	err := db.c.QueryRow("SELECT Filename FROM Photos WHERE PhotoID = ?", photoID).Scan(&filename)
	if err == sql.ErrNoRows {
		return "", errorDefinition.ErrPhotoNotFound
	} else if err != nil {
		return "", err
	}
	return filename, err
}

// TODO Comment
func (db *appdbimpl) GetLikes(photoID string, offset int, userRequesting string) ([]structs.Like, error) {
	var likes []structs.Like

	rows, err := db.c.Query("SELECT UserID FROM Likes WHERE PhotoID=? LIMIT 30 OFFSET ?", photoID, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var like structs.Like
		err = rows.Scan(&like.UserID)
		if err != nil {
			return nil, err
		}

		res, err := db.IsBanned(like.UserID, userRequesting)
		res2, err := db.IsBanned(userRequesting, like.UserID)

		if (res || res2) && err == nil {
			continue
		} else if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}
	if len(likes) == 0 {
		return nil, sql.ErrNoRows
	}
	return likes, nil
}

// TODO Comment
func (db *appdbimpl) GetComments(photoID string, offset int, userRequesting string) ([]structs.Comment, error) {
	var comments []structs.Comment

	rows, err := db.c.Query("SELECT CommentID,UserID,Text FROM Comments WHERE PhotoID=? LIMIT 30 OFFSET ?", photoID, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment structs.Comment
		err = rows.Scan(&comment.CommentID, &comment.UserID, &comment.Text)
		if err != nil {
			return nil, err
		}
		res, err := db.IsBanned(comment.UserID, userRequesting)
		res2, err := db.IsBanned(userRequesting, comment.UserID)

		if (res || res2) && err == nil {
			continue
		} else if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if len(comments) == 0 {
		return nil, sql.ErrNoRows
	}
	return comments, nil
}
