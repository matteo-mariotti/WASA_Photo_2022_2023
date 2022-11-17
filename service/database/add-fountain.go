package database

import "WASA_Photo/service/structs"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) InsertFountain(fountain structs.Fountain) error {
	_, err := db.c.Exec("INSERT INTO WASA_Photo (id, latitude, longitude, status) VALUES (?, ?,?,?)", fountain.ID, fountain.Latitude, fountain.Longitude, fountain.Status)
	return err
}
