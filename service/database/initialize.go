package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) Initialize() error {
	sql := `
			DROP TABLE IF EXISTS WASA_Photo;
			CREATE TABLE WASA_Photo(id INTEGER PRIMARY KEY, latitude FLOAT, longitude FLOAT, status STRING); 
			INSERT INTO WASA_Photo (id, latitude, longitude, status) VALUES (1, 1.0, 1.0, "good");
			INSERT INTO WASA_Photo (id, latitude, longitude, status) VALUES (2, 5.0, 6.0, "good");
			INSERT INTO WASA_Photo (id, latitude, longitude, status) VALUES (3, 7.0, 3.0, "faulty");`
	_, err := db.c.Exec(sql)
	return err
}
