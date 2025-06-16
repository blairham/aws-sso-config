package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

// Mock implementation of SSOOIDCClient interface for testing
type MockSSOOIDCClient struct {
	mock.Mock
}

func (m *MockSSOOIDCClient) RegisterClient(
	ctx context.Context,
	params *ssooidc.RegisterClientInput,
	optFns ...func(*ssooidc.Options),
) (*ssooidc.RegisterClientOutput, error) {
	args := m.Called(ctx, params, optFns)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*ssooidc.RegisterClientOutput), args.Error(1)
}

func (m *MockSSOOIDCClient) StartDeviceAuthorization(
	ctx context.Context,
	params *ssooidc.StartDeviceAuthorizationInput,
	optFns ...func(*ssooidc.Options),
) (*ssooidc.StartDeviceAuthorizationOutput, error) {
	args := m.Called(ctx, params, optFns)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*ssooidc.StartDeviceAuthorizationOutput), args.Error(1)
}

func (m *MockSSOOIDCClient) CreateToken(ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error) {
	args := m.Called(ctx, params, optFns)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*ssooidc.CreateTokenOutput), args.Error(1)
}

func TestAWSProviderGenerateTokenSuccess(t *testing.T) {
	// Create mocks and provider
	mockClient := new(MockSSOOIDCClient)

	// Success response mocks
	mockRegister := &ssooidc.RegisterClientOutput{
		ClientId:     aws.String("test-client-id"),
		ClientSecret: aws.String("test-client-secret"),
	}
	mockDeviceAuth := &ssooidc.StartDeviceAuthorizationOutput{
		DeviceCode:              aws.String("test-device-code"),
		VerificationUriComplete: aws.String("https://test-verification-uri.com"),
	}

	// Mock the client method calls
	mockClient.On("RegisterClient", mock.Anything, mock.Anything, mock.Anything).Return(mockRegister, nil)
	mockClient.On("StartDeviceAuthorization", mock.Anything, mock.Anything, mock.Anything).Return(mockDeviceAuth, nil)

	// Count browser opens
	browserOpened := false

	// Setup test provider
	provider := &AWSProvider{
		SSOOIDCClient: mockClient,
		BrowserOpener: func(url string) error {
			browserOpened = true
			assert.Equal(t, "https://test-verification-uri.com", url)
			return nil
		},
		TokenPoller: func(client SSOOIDCClient, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string {
			// Verify params passed to token poller are correct
			assert.Equal(t, mockRegister, register)
			assert.Equal(t, mockDeviceAuth, deviceAuth)
			token := "test-access-token"
			return &token
		},
	}

	// Setup test config
	appCfg := &appconfig.Config{
		SSO: appconfig.SSOConfig{
			StartURL: "https://test-sso-url.com",
			Region:   "us-west-2",
		},
	}

	// Call the method under test
	token := provider.GenerateToken(appCfg)

	// Verify results
	assert.NotNil(t, token)
	assert.Equal(t, "test-access-token", *token)
	assert.True(t, browserOpened, "Browser should have been opened")

	// Verify mock calls
	mockClient.AssertExpectations(t)
}

func TestAWSProviderGenerateTokenRegisterFailure(t *testing.T) {
	// Create mocks
	mockClient := new(MockSSOOIDCClient)

	// Mock register client failure
	mockClient.On("RegisterClient", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("register client error"))

	// Setup provider with minimal mocks
	provider := &AWSProvider{
		SSOOIDCClient: mockClient,
		BrowserOpener: func(url string) error {
			t.Fail() // Should not be called
			return nil
		},
		TokenPoller: func(client SSOOIDCClient, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string {
			t.Fail() // Should not be called
			return nil
		},
	}

	// Setup test config
	appCfg := &appconfig.Config{
		SSO: appconfig.SSOConfig{
			StartURL: "https://test-sso-url.com",
			Region:   "us-west-2",
		},
	}

	// Call the method under test
	token := provider.GenerateToken(appCfg)

	// Verify results
	assert.Nil(t, token, "Token should be nil when register client fails")

	// Verify mock calls
	mockClient.AssertExpectations(t)
}

func TestAWSProviderGenerateTokenDeviceAuthFailure(t *testing.T) {
	// Create mocks
	mockClient := new(MockSSOOIDCClient)

	// Success register response
	mockRegister := &ssooidc.RegisterClientOutput{
		ClientId:     aws.String("test-client-id"),
		ClientSecret: aws.String("test-client-secret"),
	}

	// Mock success register but failure in device auth
	mockClient.On("RegisterClient", mock.Anything, mock.Anything, mock.Anything).
		Return(mockRegister, nil)
	mockClient.On("StartDeviceAuthorization", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("device auth error"))

	// Setup provider
	provider := &AWSProvider{
		SSOOIDCClient: mockClient,
		BrowserOpener: func(url string) error {
			t.Fail() // Should not be called
			return nil
		},
		TokenPoller: func(client SSOOIDCClient, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string {
			t.Fail() // Should not be called
			return nil
		},
	}

	// Setup test config
	appCfg := &appconfig.Config{
		SSO: appconfig.SSOConfig{
			StartURL: "https://test-sso-url.com",
			Region:   "us-west-2",
		},
	}

	// Call the method under test
	token := provider.GenerateToken(appCfg)

	// Verify results
	assert.Nil(t, token, "Token should be nil when device auth fails")

	// Verify mock calls
	mockClient.AssertExpectations(t)
}

func TestAWSProviderGenerateTokenBrowserOpenFailure(t *testing.T) {
	// Create mocks
	mockClient := new(MockSSOOIDCClient)

	// Success response mocks
	mockRegister := &ssooidc.RegisterClientOutput{
		ClientId:     aws.String("test-client-id"),
		ClientSecret: aws.String("test-client-secret"),
	}
	mockDeviceAuth := &ssooidc.StartDeviceAuthorizationOutput{
		DeviceCode:              aws.String("test-device-code"),
		VerificationUriComplete: aws.String("https://test-verification-uri.com"),
	}

	// Mock successful API calls
	mockClient.On("RegisterClient", mock.Anything, mock.Anything, mock.Anything).
		Return(mockRegister, nil)
	mockClient.On("StartDeviceAuthorization", mock.Anything, mock.Anything, mock.Anything).
		Return(mockDeviceAuth, nil)

	// Setup provider with browser open failure
	provider := &AWSProvider{
		SSOOIDCClient: mockClient,
		BrowserOpener: func(url string) error {
			return errors.New("browser open error")
		},
		TokenPoller: func(client SSOOIDCClient, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string {
			// Should still be called even if browser fails
			token := "test-access-token"
			return &token
		},
	}

	// Setup test config
	appCfg := &appconfig.Config{
		SSO: appconfig.SSOConfig{
			StartURL: "https://test-sso-url.com",
			Region:   "us-west-2",
		},
	}

	// Call the method under test
	token := provider.GenerateToken(appCfg)

	// Verify results - should continue even if browser open fails
	assert.NotNil(t, token)
	assert.Equal(t, "test-access-token", *token)

	// Verify mock calls
	mockClient.AssertExpectations(t)
}

// TestNewDefaultAWSProvider tests the NewDefaultAWSProvider function
func TestNewDefaultAWSProvider(t *testing.T) {
	provider := NewDefaultAWSProvider()

	assert.NotNil(t, provider, "Provider should not be nil")
	assert.NotNil(t, provider.SSOOIDCClient, "SSOOIDCClient should not be nil")
	assert.NotNil(t, provider.BrowserOpener, "BrowserOpener should not be nil")
	assert.NotNil(t, provider.TokenPoller, "TokenPoller should not be nil")
	assert.NotNil(t, provider.Cfg, "Cfg should not be nil")
}
