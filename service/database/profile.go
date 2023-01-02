package database

import (
	"WASA_Photo/service/structs"
	"database/sql"
	"strconv"
)

// GetFollowerNumber is a function that returns the number of followers of a user
func (db *appdbimpl) GetFollowerNumber(userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE UserB=? AND (Followers.UserA,?) NOT IN (SELECT * FROM Bans) AND (?,Followers.UserA) NOT IN (SELECT * FROM Bans)", userID, userID, userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}

// GetFollowingNumber is a function that returns the number of following of a user
func (db *appdbimpl) GetFollowingNumber(userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE UserA=? AND (UserA,Followers.UserB) NOT IN (SELECT * FROM Bans) AND (Followers.UserB,UserA) NOT IN (SELECT * FROM Bans)", userID, userID, userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}

// GetPhotosNumber is a function that returns the number of photos of a user
func (db *appdbimpl) GetPhotosNumber(userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Photos WHERE Owner=?", userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}

// GetPhoto is a function that returns the photo with the given ID
func (db *appdbimpl) GetPhotos(userID string, reqUser string, offset int) ([]structs.Photo, error) {
	type photoPartialInfo struct {
		PhotoID int
		Owner   string
		Date    string
	}

	var photos []photoPartialInfo
	rows, err := db.c.Query("SELECT PhotoID, Owner, Date FROM Photos WHERE Owner=? ORDER BY PhotoID DESC LIMIT 30 OFFSET ?", userID, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var photo photoPartialInfo
		err = rows.Scan(&photo.PhotoID, &photo.Owner, &photo.Date)
		if err != nil {
			return nil, err
		}
		// Get UserName instead of tokenID
		user := db.c.QueryRow("SELECT UserName FROM Users WHERE UserID=?", photo.Owner)
		err = user.Scan(&photo.Owner)
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
		likes, err := db.getLikesNumber(photos[i].PhotoID, reqUser)
		if err != nil {
			return nil, err
		}
		// Get number of comments
		comments, err := db.getCommentsNumber(photos[i].PhotoID, reqUser)
		if err != nil {
			return nil, err
		}
		userLike, err := db.HasLiked(strconv.Itoa(photos[i].PhotoID), reqUser)
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
			LoggedLike:     userLike,
		})

	}
	return photosInfo, nil
}

// GetLikesNumber is a function that returns the number of likes of a photo
func (db *appdbimpl) getLikesNumber(photoID int, userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Likes, Photos WHERE Likes.PhotoID=? AND Photos.PhotoID=Likes.PhotoID AND (UserID,?) NOT IN (SELECT * FROM Bans) AND (?,UserID) NOT IN (SELECT * FROM Bans) AND (Photos.Owner, UserID) NOT IN (SELECT * FROM Bans) AND (UserID, Photos.Owner) NOT IN (SELECT * FROM Bans)", photoID, userID, userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}

// GetCommentsNumber is a function that returns the number of comments of a photo
func (db *appdbimpl) getCommentsNumber(photoID int, userID string) (int, error) {
	var count int
	row := db.c.QueryRow("SELECT COUNT(*) FROM Comments,Photos WHERE Comments.PhotoID=? AND Photos.PhotoID=Comments.PhotoID AND (UserID,?) NOT IN (SELECT * FROM Bans) AND (?,UserID) NOT IN (SELECT * FROM Bans) AND (Photos.Owner, UserID) NOT IN (SELECT * FROM Bans) AND (UserID, Photos.Owner) NOT IN (SELECT * FROM Bans)", photoID, userID, userID)
	// Note that we are not checking for sql.ErrNoRows here because we are using count(*) and it will always return a row
	err := row.Scan(&count)
	return count, err
}
