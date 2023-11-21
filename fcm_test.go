package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFMCService(t *testing.T) {
	// Mocking
	bgCtx := context.Background()

	credentialJSON := []byte(`{
        "type": "service_account",
        "project_id": "test-project",
        "private_key_id": "test-key-id",
        "private_key": "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
        "client_email": "test@test-project.iam.gserviceaccount.com",
        "client_id": "test-client-id",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com"
    }`)

	_, err := NewFMCService(nil, bgCtx, credentialJSON)
	assert.NoError(t, err)
}
