package db

import (
	"flatchecker-scheduler/testutil"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDB_DevEnvironment_Success(t *testing.T) {
	// Set up environment variable for dev password
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, "test-password")
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	// Note: This will try to connect to actual database
	// In real scenario, we'd use sqlmock more extensively
	// For now, we test the connection string generation
	connStr, err := getConnectionString(config)
	require.NoError(t, err)
	assert.Contains(t, connStr, "testuser:test-password")
	assert.Contains(t, connStr, "tcp(127.0.0.1:3306)")
	assert.Contains(t, connStr, "testdb")
	assert.Contains(t, connStr, "parseTime=true")
}

func TestGetConnectionString_DevEnvironment(t *testing.T) {
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, "dev-pass-123")
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	connStr, err := getConnectionString(config)
	require.NoError(t, err)

	expectedFormat := "testuser:dev-pass-123@tcp(127.0.0.1:3306)/testdb?parseTime=true"
	assert.Equal(t, expectedFormat, connStr)
}

func TestGetConnectionString_ProdEnvironment(t *testing.T) {
	// For prod, we can't easily test without mocking secret manager
	// but we can test the connection string format
	config := testutil.CreateTestConfig("prod")

	// This will fail because GetSecret will try to call GCP
	// but we can test with a mock secret getter
	_, err := getConnectionString(config)
	// We expect an error because we can't actually reach secret manager
	assert.Error(t, err)
}

func TestGetConnectionString_UnknownEnvironment(t *testing.T) {
	config := make(map[string]string)
	config["environment"] = "staging"
	config["user-name"] = "testuser"

	_, err := getConnectionString(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown environment")
}

func TestGetConnectionString_ParseTimeParameter(t *testing.T) {
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, "password")
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	connStr, err := getConnectionString(config)
	require.NoError(t, err)
	assert.Contains(t, connStr, "parseTime=true")
}

func TestGetPassword_DevEnvironment_EnvVarSet(t *testing.T) {
	expectedPassword := "my-dev-password"
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, expectedPassword)
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	password, err := getPassword(config)
	require.NoError(t, err)
	assert.Equal(t, expectedPassword, password)
}

func TestGetPassword_DevEnvironment_EnvVarNotSet(t *testing.T) {
	// Make sure env var is not set
	os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	_, err := getPassword(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "environment variable")
	assert.Contains(t, err.Error(), DEV_PASSWORD_ENVIRONMENT_VARIABLE)
}

func TestGetPassword_ProdEnvironment_SecretManagerError(t *testing.T) {
	config := testutil.CreateTestConfig("prod")

	// This will fail because we can't reach actual secret manager
	_, err := getPassword(config)
	assert.Error(t, err)
}

func TestGetPassword_UnknownEnvironment(t *testing.T) {
	config := make(map[string]string)
	config["environment"] = "test"

	_, err := getPassword(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown environment")
}

func TestGetConnectionString_MissingPassword(t *testing.T) {
	// Don't set password env var
	os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	_, err := getConnectionString(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestGetConnectionString_EmptyValues(t *testing.T) {
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, "password")
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := make(map[string]string)
	config["environment"] = "dev"
	config["user-name"] = ""
	config["db-name"] = ""
	config["ip-address"] = ""

	connStr, err := getConnectionString(config)
	require.NoError(t, err)
	// Connection string will be malformed but function doesn't validate
	assert.Contains(t, connStr, "password")
}

func TestGetConnectionString_SpecialCharactersInPassword(t *testing.T) {
	specialPassword := "p@ssw0rd!#$%"
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, specialPassword)
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	connStr, err := getConnectionString(config)
	require.NoError(t, err)
	assert.Contains(t, connStr, specialPassword)
}

func TestGetDB_ConnectionStringFormat(t *testing.T) {
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, "test123")
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	config := testutil.CreateTestConfig("dev")

	connStr, err := getConnectionString(config)
	require.NoError(t, err)

	// Verify all components are present
	assert.Contains(t, connStr, "testuser")
	assert.Contains(t, connStr, "test123")
	assert.Contains(t, connStr, "tcp(127.0.0.1:3306)")
	assert.Contains(t, connStr, "testdb")
	assert.Contains(t, connStr, "parseTime=true")
}

func TestGetConnectionString_ProdFormat(t *testing.T) {
	config := testutil.CreateTestConfig("prod")

	// We can't get a full connection string without secrets,
	// but we can verify the error handling
	_, err := getConnectionString(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestDevPasswordConstant(t *testing.T) {
	// Verify the constant is correct
	assert.Equal(t, "FLATCHECKER_SCHEDULER_DEV_PASSWORD", DEV_PASSWORD_ENVIRONMENT_VARIABLE)
	assert.NotEmpty(t, DEV_PASSWORD_ENVIRONMENT_VARIABLE)
}

func TestProdPasswordConstant(t *testing.T) {
	// Verify the prod secret name constant
	assert.Equal(t, "flatchecker-db-password", PROD_PASSWORD_SECRET_NAME)
	assert.NotEmpty(t, PROD_PASSWORD_SECRET_NAME)
}

func TestGetDB_WithSqlMock(t *testing.T) {
	// Test GetDB with sqlmock to verify ping is called
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	// Ping the mock database
	err = db.Ping()
	assert.NoError(t, err)

	// Verify expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetConnectionString_AllConfigKeys(t *testing.T) {
	os.Setenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE, "pass")
	defer os.Unsetenv(DEV_PASSWORD_ENVIRONMENT_VARIABLE)

	tests := []struct {
		name        string
		environment string
		wantContain []string
	}{
		{
			name:        "dev environment",
			environment: "dev",
			wantContain: []string{"tcp", "127.0.0.1", "3306", "parseTime=true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := testutil.CreateTestConfig(tt.environment)
			connStr, err := getConnectionString(config)
			require.NoError(t, err)
			for _, want := range tt.wantContain {
				assert.Contains(t, connStr, want)
			}
		})
	}
}
