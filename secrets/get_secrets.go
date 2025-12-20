package secrets

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

const GCP_PROJECT = "flatchecker"

// SecretGetter is an interface for retrieving secrets
type SecretGetter interface {
	GetSecret(name string) (string, error)
}

// GCPSecretGetter implements SecretGetter using Google Cloud Secret Manager
type GCPSecretGetter struct{}

// NewSecretGetter creates a new SecretGetter
func NewSecretGetter() SecretGetter {
	return &GCPSecretGetter{}
}

// GetSecret retrieves a secret from Google Cloud Secret Manager
func (g *GCPSecretGetter) GetSecret(name string) (string, error) {
	secretVersionName := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", GCP_PROJECT, name)

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretVersionName,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}

// GetSecret is a convenience function that uses the default GCP implementation
func GetSecret(name string) (string, error) {
	getter := NewSecretGetter()
	return getter.GetSecret(name)
}
