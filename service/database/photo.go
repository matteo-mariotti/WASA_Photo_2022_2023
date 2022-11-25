package database

// FollowUser is a function that adds userB to the list of users that userA is following
func (db *appdbimpl) UploadPhoto(owner string, filename string) error {
	_, err := db.c.Exec("INSERT INTO Photos (Owner, Filename) VALUES (?, ?)", owner, filename)
	return err
}

// FollowUser is a function that removes userB from the list of users userA is following
func (db *appdbimpl) DeletePhoto(photoID string) (string, error) {
	var fileIdentifier string
	err := db.c.QueryRow("SELECT Filename FROM Photos WHERE PhotoID = ?", photoID).Scan(&fileIdentifier)
	if err != nil {
		return "", err
	}
	_, err = db.c.Exec("DELETE FROM Photos WHERE PhotoID=?", photoID)
	return fileIdentifier, err
}

func (db *appdbimpl) GetPhotoOwner(photoID string) (string, error) {
	var owner string
	err := db.c.QueryRow("SELECT Owner FROM Photos WHERE PhotoID = ?", photoID).Scan(&owner)
	return owner, err
}
