package db

import (
	"database/sql"
	"fmt"
)

func GetDB(config map[string]string) (*sql.DB, error){
	connectionString:=fmt.Sprintf("%s:%s@%s/%s",config["user-name"], config["password"], config["ip-address"], config["db-name"])

	db, err := sql.Open("mysql", connectionString)
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