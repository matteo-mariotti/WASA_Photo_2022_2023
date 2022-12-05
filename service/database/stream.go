package database

import (
	"WASA_Photo/service/structs"
	"database/sql"
	"strconv"
)

// GetFollowingPhotosChrono is a function that returns the photos of the users that the user is following in reverse chronological order
func (db *appdbimpl) GetFollowingPhotosChrono(following []string, offset int) ([]structs.Photo, error) {
	type photoPartialInfo struct {
		PhotoID int
		Owner   string
		Date    string
	}

	// Create a string with the following users
	var followingString string
	for i := range following {
		if i == 0 {
			followingString = "'" + following[i] + "'"
		} else {
			followingString += "," + "'" + following[i] + "'"
		}
	}

	var photos []photoPartialInfo
	query := "SELECT PhotoID, Owner, Date FROM Photos WHERE Owner IN (" + followingString + ") ORDER BY Date DESC LIMIT 30 OFFSET " + strconv.Itoa(offset)
	rows, err := db.c.Query(query)
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