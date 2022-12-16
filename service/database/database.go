/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"WASA_Photo/service/structs"
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Transaction
	StartTransaction() error
	Commit() error
	Rollback() error

	// Check token
	ValidToken(token string) (bool, error)

	// Login functions
	RegisterUser(uuid string, username string) error // Insert a new user in the database
	LoginUser(username string) (string, error)       // Login a user

	// User functions
	UserAvailable(username string) (bool, error)          // Check if a username is available
	ChangeUsername(usedname string, newname string) error // Change the username of a user

	// Ban functions
	BanUser(userA string, userB string) error          // Block a user
	UnbanUser(userA string, userB string) error        // Unblock a user
	IsBanned(userA string, userB string) (bool, error) // Check if userB has been blocked by userA

	// Follow functions
	FollowUser(userA string, userB string) error          // Follow a user
	UnfollowUser(userA string, userB string) error        // Unfollow a user
	IsFollowing(userA string, userB string) (bool, error) // Check if userA is following userB

	// Photo functions
	UploadPhoto(owner string, filename string) error                                          // Upload a photo
	DeletePhoto(photoID string) (string, error)                                               // Delete a photo
	GetPhotoOwner(photoID string) (string, error)                                             // Get the owner of a photo
	GetPhoto(photoID string) (string, error)                                                  // Get the filename of a photo giver its ID
	GetLikes(photoID string, offset int, requestingUser string) ([]structs.Like, error)       // Get the list of users that liked a photo
	GetComments(photoID string, offset int, requestingUser string) ([]structs.Comment, error) // Get the list of comments of a photo

	// Comment functions
	Comment(photoID string, userID string, text string) error        // Comment a photo
	Uncomment(photoID string, userID string, commentID string) error // Uncomment a photo

	// Like functions
	Like(photoID string, userID string) error             // Like a photo
	Unlike(photoID string, userID string) error           // Unlike a photo
	HasLiked(photoID string, userID string) (bool, error) // Check if userID has liked a photo with the given ID

	// Profile functions
	GetName(userID string) (string, error)                                        // Get the username from the user ID
	GetFollowerNumber(userID string) (int, error)                                 // Get the number of followers of a user
	GetFollowingNumber(userID string) (int, error)                                // Get the number of users a user is following
	GetPhotosNumber(userID string) (int, error)                                   // Get the number of photos a user has uploaded
	GetPhotos(userID string, reqUser string, offset int) ([]structs.Photo, error) // Get the photos of a user

	// Search functions
	GetUsers(start string, offset int) ([]string, error) // Get the list of users that start with the given string

	// Stream
	GetFollowing(userID string) ([]string, error)                                     // Get the list of users that a user is following
	GetFollowingPhotosChrono(following []string, offset int) ([]structs.Photo, error) // Get the photos of the users that a user is following

	GetVersion() (string, error) // Get database version
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `PRAGMA foreign_keys = ON;

		BEGIN TRANSACTION;
		
		create table Users
		(
			UserID   TEXT not null
				constraint Users_pk
					primary key,
			UserName TEXT
		);
		
		create table Bans
		(
			UserA TEXT
				constraint Bans_Users_UserID_fk
					references Users,
			UserB TEXT
				constraint Bans_Users__fk2
					references Users,
			constraint Bans_pk
				primary key (UserA, UserB)
		);
		
		create table Followers
		(
			UserA TEXT
				constraint Followers___fk1
					references Users
					on update cascade on delete cascade,
			UserB TEXT
				constraint Followers___fk2
					references Users,
			constraint Followers_pk
				primary key (UserA, UserB)
		);
		
		create table Photos
		(
			PhotoID  INTEGER
				constraint Photos_pk
					primary key autoincrement,
			Owner    TEXT
				constraint Photos_Users_UserID_fk
					references Users
					on update cascade on delete cascade,
			Filename TEXT not null
		);
		
		create table Comments
		(
			CommentID INTEGER
				constraint Comments_pk
					primary key autoincrement,
			PhotoID   INTEGER not null
				references Photos
					on update cascade on delete cascade,
			UserID    TEXT
				constraint Comments_Users
					references Users
					on update cascade on delete cascade,
			Text      TEXT    not null
		);
		
		create table Likes
		(
			UserID  TEXT
				constraint Likes_Users_UserID_fk
					references Users
					on update cascade on delete cascade,
			PhotoID INTEGER
				constraint Likes_Photos_PhotoID_fk
					references Photos
					on update cascade on delete cascade,
			constraint Likes_pk
				primary key (UserID, PhotoID)
		);
		
		create table example_table
		(
			id   INTEGER not null
				primary key,
			name TEXT
		);
				
		
		COMMIT;
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	} else {
		_, err := db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
