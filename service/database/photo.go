package database

import (
	"WASA_Photo/service/errorDefinition"
	"database/sql"
)

// TODO Comment
func (db *appdbimpl) UploadPhoto(owner string, filename string) error {
	result, err := db.c.Exec("INSERT INTO Photos (Owner, Filename) VALUES (?, ?)", owner, filename)

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
