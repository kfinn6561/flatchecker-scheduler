package db

import (
	"database/sql"
)

func GetDB() (*sql.DB, error){
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}