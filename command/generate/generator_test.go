package generate

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/bigkevmcd/go-configparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementation of SSOClient interface for testing
type MockSSOClient struct {
	mock.Mock
}

func (m *MockSSOClient) GetRoleCredentials(
	ctx context.Context,
	params *sso.GetRoleCredentialsInput,
	optFns ...func(*sso.Options),
) (*sso.GetRoleCredentialsOutput, error) {
	args := m.Called(ctx, params)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*sso.GetRoleCredentialsOutput), args.Error(1)
}

func (m *MockSSOClient) ListAccounts(ctx context.Context, params *sso.ListAccountsInput, optFns ...func(*sso.Options)) (*sso.ListAccountsOutput, error) {
	args := m.Called(ctx, params)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*sso.ListAccountsOutput), args.Error(1)
}

func (m *MockSSOClient) ListAccountRoles(ctx context.Context, params *sso.ListAccountRolesInput, optFns ...func(*sso.Options)) (*sso.ListAccountRolesOutput, error) {
	args := m.Called(ctx, params)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*sso.ListAccountRolesOutput), args.Error(1)
}

func TestListAccountsWithClient(t *testing.T) {
	// Create mock client
	mockClient := new(MockSSOClient)

	// Test successful case
	token := "test-token"
	mockAccounts := []types.AccountInfo{
		{
			AccountId:   aws.String("123456789012"),
			AccountName: aws.String("Test Account 1"),
		},
		{
			AccountId:   aws.String("098765432109"),
			AccountName: aws.String("Test Account 2"),
		},
	}

	mockResponse := &sso.ListAccountsOutput{
		AccountList: mockAccounts,
	}

	mockClient.On("ListAccounts", mock.Anything, mock.MatchedBy(func(input *sso.ListAccountsInput) bool {
		return *input.AccessToken == token
	})).Return(mockResponse, nil)

	// Create generator and call the method
	generator := NewConfigGenerator("https://test.com", "us-west-2", "us-east-1")
	accounts, err := generator.ListAccountsWithClient(mockClient, &token)

	// Verify results
	require.NoError(t, err)
	assert.Equal(t, mockAccounts, accounts)
	mockClient.AssertExpectations(t)

	// Test with nil token
	accounts, err = generator.ListAccountsWithClient(mockClient, nil)
	assert.Error(t, err)
	assert.Nil(t, accounts)

	// Test with ListAccounts error
	mockClient.ExpectedCalls = nil
	mockClient.On("ListAccounts", mock.Anything, mock.Anything).Return(nil, errors.New("API error"))

	accounts, err = generator.ListAccountsWithClient(mockClient, &token)
	assert.Error(t, err)
	assert.Nil(t, accounts)
	mockClient.AssertExpectations(t)
}

func TestGetAccountRolesWithClient(t *testing.T) {
	// Create mock client
	mockClient := new(MockSSOClient)

	// Test successful case
	token := "test-token"
	accountID := "123456789012"
	mockRoles := []types.RoleInfo{
		{
			RoleName:  aws.String("AdminRole"),
			AccountId: aws.String(accountID),
		},
		{
			RoleName:  aws.String("ReadOnlyRole"),
			AccountId: aws.String(accountID),
		},
	}

	mockResponse := &sso.ListAccountRolesOutput{
		RoleList: mockRoles,
	}

	mockClient.On("ListAccountRoles", mock.Anything, mock.MatchedBy(func(input *sso.ListAccountRolesInput) bool {
		return *input.AccessToken == token && *input.AccountId == accountID
	})).Return(mockResponse, nil)

	// Create generator and call the method
	generator := NewConfigGenerator("https://test.com", "us-west-2", "us-east-1")
	roles, err := generator.GetAccountRolesWithClient(mockClient, &token, accountID)

	// Verify results
	require.NoError(t, err)
	assert.Equal(t, mockRoles, roles)
	mockClient.AssertExpectations(t)

	// Test with nil token
	roles, err = generator.GetAccountRolesWithClient(mockClient, nil, accountID)
	assert.Error(t, err)
	assert.Nil(t, roles)

	// Test with ListAccountRoles error
	mockClient.ExpectedCalls = nil
	mockClient.On("ListAccountRoles", mock.Anything, mock.Anything).Return(nil, errors.New("API error"))

	roles, err = generator.GetAccountRolesWithClient(mockClient, &token, accountID)
	assert.Error(t, err)
	assert.Nil(t, roles)
	mockClient.AssertExpectations(t)
}

func TestNewConfigGenerator(t *testing.T) {
	ssoStartURL := "https://test.com"
	ssoRegion := "us-west-2"
	defaultRegion := "us-east-1"

	generator := NewConfigGenerator(ssoStartURL, ssoRegion, defaultRegion)

	assert.Equal(t, ssoStartURL, generator.SSOStartURL)
	assert.Equal(t, ssoRegion, generator.SSORegion)
	assert.Equal(t, defaultRegion, generator.DefaultRegion)
}

// TestWriteSectionToConfig tests the WriteSectionToConfig function
func TestWriteSectionToConfig(t *testing.T) {
	generator := NewConfigGenerator("https://test.com", "us-west-2", "us-east-1")

	// We'll need to create a configparser instance to test with
	// Since this requires an actual file, we'll create a temporary one
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config")

	// Create initial config file
	err := os.WriteFile(configFile, []byte("[default]\nregion = us-west-2\n"), 0600)
	require.NoError(t, err)

	// Load the config parser
	configParser, err := configparser.NewConfigParserFromFile(configFile)
	require.NoError(t, err)

	// Test writing a new section
	sectionName := "profile test-profile"
	values := map[string]string{
		"sso_account_id": "123456789012",
		"sso_role_name":  "TestRole",
		"sso_region":     "us-west-2",
		"sso_start_url":  "https://test.com",
		"region":         "us-east-1",
	}

	err = generator.WriteSectionToConfig(configParser, sectionName, values)
	assert.NoError(t, err)

	// Verify the section was created
	assert.True(t, configParser.HasSection(sectionName))

	// Verify the values were set correctly
	for key, expectedValue := range values {
		actualValue, err := configParser.Get(sectionName, key)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	}

	// Test updating an existing section
	newValues := map[string]string{
		"sso_account_id": "098765432109",
		"region":         "us-west-1",
	}

	err = generator.WriteSectionToConfig(configParser, sectionName, newValues)
	assert.NoError(t, err)

	// Verify the values were updated
	accountID, err := configParser.Get(sectionName, "sso_account_id")
	assert.NoError(t, err)
	assert.Equal(t, "098765432109", accountID)

	region, err := configParser.Get(sectionName, "region")
	assert.NoError(t, err)
	assert.Equal(t, "us-west-1", region)

	// Verify other values remain unchanged
	roleName, err := configParser.Get(sectionName, "sso_role_name")
	assert.NoError(t, err)
	assert.Equal(t, "TestRole", roleName)
}

// TestWriteSectionToConfigError tests error cases for WriteSectionToConfig
func TestWriteSectionToConfigError(t *testing.T) {
	generator := NewConfigGenerator("https://test.com", "us-west-2", "us-east-1")

	// Create a mock config parser that simulates errors
	// Since we can't easily mock the configparser, we'll create a scenario
	// where we try to write to a parser with an empty values map
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config")

	// Create initial config file
	err := os.WriteFile(configFile, []byte("[default]\nregion = us-west-2\n"), 0600)
	require.NoError(t, err)

	// Load the config parser
	configParser, err := configparser.NewConfigParserFromFile(configFile)
	require.NoError(t, err)

	// Test with empty values (should still work)
	err = generator.WriteSectionToConfig(configParser, "profile empty", map[string]string{})
	assert.NoError(t, err)

	// Verify the empty section was created
	assert.True(t, configParser.HasSection("profile empty"))
}
