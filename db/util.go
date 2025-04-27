package db

import (
	"os"
	"path/filepath"
)

const SQL_FOLDER = "sql"

func readSqlFile(sqlName string) (string, error) {
	path := filepath.Join(SQL_FOLDER, sqlName)
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
