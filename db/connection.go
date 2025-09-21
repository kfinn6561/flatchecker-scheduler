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

	connectionString, err := getConnectionString(config)
	if err != nil {
		return nil, fmt.Errorf("error getting connection string: %v", err)
	}

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

func getConnectionString(config map[string]string) (string, error) {
	password, err := getPassword(config)
	if err != nil {
		return "", fmt.Errorf("error getting password: %v", err)
	}
	switch config["environment"] {
	case "dev":
		return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config["user-name"], password, config["ip-address"], config["db-name"]), nil
	case "prod":
		return fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true", config["user-name"], password, config["connection-name"], config["db-name"]), nil
	default:
		return "", fmt.Errorf("unknown environment: %s", config["environment"])
	}
}

func getPassword(config map[string]string) (string, error) {
	switch config["environment"] {
	case "dev":
		password, ok := os.LookupEnv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)
		if ok {
			return password, nil
		} else {
			return "", fmt.Errorf("environment variable %s not set", DEV_PASSWORD_ENVIRONMENT_VARIABLE)
		}
	case "prod":
		return secrets.GetSecret(PROD_PASSWORD_SECRET_NAME)
	default:
		return "", fmt.Errorf("unknown environment: %s", config["environment"])
	}
}
