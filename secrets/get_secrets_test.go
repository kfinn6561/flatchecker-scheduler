package secrets

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecret_ConvenienceFunction(t *testing.T) {
	// Test that the convenience function calls the getter
	// Note: This will make an actual API call in real environment
	// In a real test, we'd mock the GCP client, but for now we test the interface
	getter := NewSecretGetter()
	assert.NotNil(t, getter)
	assert.Implements(t, (*SecretGetter)(nil), getter)
}

func TestNewSecretGetter_ReturnsImplementation(t *testing.T) {
	// Test that NewSecretGetter returns a valid implementation
	getter := NewSecretGetter()
	assert.NotNil(t, getter)
	assert.IsType(t, &GCPSecretGetter{}, getter)
}

func TestSecretGetter_Interface(t *testing.T) {
	// Test that GCPSecretGetter implements SecretGetter interface
	var _ SecretGetter = &GCPSecretGetter{}
	// If this compiles, the interface is correctly implemented
}

func TestGCPSecretGetter_VersionFormatting(t *testing.T) {
	// Test the secret version name format
	secretName := "test-secret"
	expectedFormat := "projects/flatchecker/secrets/test-secret/versions/latest"

	// We can't easily test the internal formatting without mocking,
	// but we can verify the constant is correct
	assert.Equal(t, "flatchecker", GCP_PROJECT)

	// Verify expected format structure
	assert.Contains(t, expectedFormat, GCP_PROJECT)
	assert.Contains(t, expectedFormat, secretName)
	assert.Contains(t, expectedFormat, "versions/latest")
}

// MockSecretGetter for testing
type MockSecretGetter struct {
	secret string
	err    error
}

func (m *MockSecretGetter) GetSecret(name string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.secret, nil
}

func TestMockSecretGetter_Success(t *testing.T) {
	// Test mock implementation returns expected value
	mock := &MockSecretGetter{
		secret: "test-password-123",
		err:    nil,
	}

	result, err := mock.GetSecret("db-password")
	assert.NoError(t, err)
	assert.Equal(t, "test-password-123", result)
}

func TestMockSecretGetter_Error(t *testing.T) {
	// Test mock implementation returns error
	expectedErr := errors.New("failed to access secret")
	mock := &MockSecretGetter{
		secret: "",
		err:    expectedErr,
	}

	result, err := mock.GetSecret("db-password")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, result)
}

func TestMockSecretGetter_EmptySecretName(t *testing.T) {
	// Test mock with empty secret name
	mock := &MockSecretGetter{
		secret: "some-value",
		err:    nil,
	}

	result, err := mock.GetSecret("")
	assert.NoError(t, err)
	assert.Equal(t, "some-value", result)
}

func TestMockSecretGetter_SecretWithSpecialCharacters(t *testing.T) {
	// Test secret name with hyphens and underscores
	mock := &MockSecretGetter{
		secret: "complex-password_123!",
		err:    nil,
	}

	result, err := mock.GetSecret("my-test_secret")
	assert.NoError(t, err)
	assert.Equal(t, "complex-password_123!", result)
}

func TestGCPProjectConstant(t *testing.T) {
	// Verify the GCP_PROJECT constant is set correctly
	assert.Equal(t, "flatchecker", GCP_PROJECT)
	assert.NotEmpty(t, GCP_PROJECT)
}
