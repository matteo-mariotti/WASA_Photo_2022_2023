package database

import (
	"WASA_Photo/service/structs"
	"database/sql"
)

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

func (db *appdbimpl) GetPhotos(userID string, offset int) ([]structs.Photo, error) {
	type photoPartialInfo struct {
		PhotoID int
		Owner   string
		Date    string
	}

	var photos []photoPartialInfo
	rows, err := db.c.Query("SELECT PhotoID, Owner, Date FROM Photos WHERE Owner=? ORDER BY PhotoID DESC LIMIT 30 OFFSET ?", userID, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var photo photoPartialInfo
		err = rows.Scan(&photo.PhotoID, &photo.Owner, &photo.Date)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	if len(photos) == 0 {
		return nil, sql.ErrNoRows
	}

	var photosInfo []structs.Photo

	// Get number of likes and comments for each photo
	for i := range photos {
		// Get number of likes
		likes, err := db.getLikesNumber(photos[i].PhotoID)
		if err != nil {
			return nil, err
		}
		// Get number of comments
		comments, err := db.getCommentsNumber(photos[i].PhotoID)
		if err != nil {
			return nil, err
		}

		// Append photo info to photosInfo
		photosInfo = append(photosInfo, structs.Photo{
			PhotoID:        photos[i].PhotoID,
			Owner:          photos[i].Owner,
			Date:           photos[i].Date,
			LikesNumber:    likes,
			CommentsNumber: comments,
		})

	}
	return photosInfo, nil
}

func (db *appdbimpl) getLikesNumber(photoID int) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PhotoID=?", photoID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}

func (db *appdbimpl) getCommentsNumber(photoID int) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PhotoID=?", photoID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}
