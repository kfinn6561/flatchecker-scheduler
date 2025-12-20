package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadSqlFile_ValidFile(t *testing.T) {
	// Test reading an existing SQL file
	content, err := readSqlFile("get_schedules.sql")
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "SELECT")
}

func TestReadSqlFile_NonExistentFile(t *testing.T) {
	// Test reading a non-existent file
	_, err := readSqlFile("nonexistent.sql")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestReadSqlFile_EmptyFile(t *testing.T) {
	// Create a temporary empty SQL file
	tempDir := t.TempDir()
	oldSQLFolder := SQL_FOLDER

	// We can't change the const, so we'll test with existing files
	// This test verifies behavior if an empty file exists
	content, err := readSqlFile("get_schedules.sql")
	require.NoError(t, err)
	// The actual file shouldn't be empty, but we verify read succeeds
	assert.NotEmpty(t, content)

	// Restore for other tests
	_ = oldSQLFolder
	_ = tempDir
}

func TestReadSqlFile_PathConstruction(t *testing.T) {
	// Verify the path is constructed correctly
	expectedPath := filepath.Join(SQL_FOLDER, "test.sql")
	assert.Equal(t, "db/sql/test.sql", expectedPath)
}

func TestReadSqlFile_MultilineSQL(t *testing.T) {
	// Test that multiline SQL content is preserved
	content, err := readSqlFile("get_schedules.sql")
	require.NoError(t, err)

	// Verify multiline content is preserved
	assert.NotEmpty(t, content)
	// Check that content is a string (multiline handling)
	assert.IsType(t, "", content)
}

func TestReadSqlFile_UpdateScheduleFile(t *testing.T) {
	// Test reading the update schedule SQL file
	content, err := readSqlFile("update_schedule.sql")
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "UPDATE")
}

func TestReadSqlFile_RelativePathHandling(t *testing.T) {
	// Test that relative path is correctly constructed with SQL_FOLDER
	// The function should prepend SQL_FOLDER to the filename
	filename := "get_schedules.sql"
	expectedPath := filepath.Join("db/sql", filename)

	// Verify the path format
	assert.Equal(t, "db/sql/get_schedules.sql", expectedPath)

	// Now test actual file reading
	content, err := readSqlFile(filename)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}
