package database

func (db *appdbimpl) StartTransaction() error {
	_, err := db.c.Exec("BEGIN TRANSACTION")
	return err
}

func (db *appdbimpl) Commit() error {
	_, err := db.c.Exec("COMMIT")
	return err
}

func (db *appdbimpl) Rollback() error {
	_, err := db.c.Exec("ROLLBACK")
	return err
}
