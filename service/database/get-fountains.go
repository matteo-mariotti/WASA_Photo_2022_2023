package database

import "WASA_Photo/service/structs"

func (db *appdbimpl) GetWASA_Photo() ([]structs.Fountain, error) {
	rows, err := db.c.Query("SELECT * FROM WASA_Photo")
	WASA_Photo := []structs.Fountain{}
	for rows.Next() {
		var id int
		var latitude float64
		var longitude float64
		var status string
		err = rows.Scan(&id, &latitude, &longitude, &status)
		WASA_Photo = append(WASA_Photo, structs.Fountain{
			ID:        id,
			Latitude:  latitude,
			Longitude: longitude,
			Status:    status,
		})
	}
	return WASA_Photo, err
}
