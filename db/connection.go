package db

import (
	"database/sql"
	"flatchecker-scheduler/secrets"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DEV_PASSWORD_ENVIRONMENT_VARIABLE = "FLATCHECKER_SCHEDULER_DEV_PASSWORD"
	PROD_PASSWORD_SECRET_NAME         = "flatchecker-db-password"
)

func GetDB(config map[string]string) (*sql.DB, error) {
	password, err := getPassword(config)
	if err != nil {
		return nil, fmt.Errorf("error getting password: %v", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config["user-name"], password, config["ip-address"], config["db-name"])

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

func getPassword(config map[string]string) (string, error) {
	if config["environment"] == "dev" {
		password, ok := os.LookupEnv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)
		if ok {
			return password, nil
		} else {
			return "", fmt.Errorf("environment variable %s not set", DEV_PASSWORD_ENVIRONMENT_VARIABLE)
		}
	} else if config["environment"] == "prod" {
		return secrets.GetSecret(PROD_PASSWORD_SECRET_NAME)
	} else {
		return "", fmt.Errorf("unknown environment: %s", config["environment"])
	}
}
