package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB(config map[string]string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config["user-name"], config["password"], config["ip-address"], config["db-name"])

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
