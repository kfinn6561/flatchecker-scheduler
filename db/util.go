package db

import (
	"os"
	"path/filepath"
)

const SQL_FOLDER = "db/sql"

func readSqlFile(sqlName string) (string, error) {
	// Try multiple paths to support both production and test execution
	paths := []string{
		filepath.Join(SQL_FOLDER, sqlName),    // Production: db/sql/file.sql
		filepath.Join("sql", sqlName),          // Test from db/: sql/file.sql
		filepath.Join("..", SQL_FOLDER, sqlName), // Test from root: ../db/sql/file.sql
	}

	var lastErr error
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		if err == nil {
			return string(contents), nil
		}
		lastErr = err
	}

	// Return the last error if none of the paths worked
	return "", lastErr
}
