package generate

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/bigkevmcd/go-configparser"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

// MockSSOClientForConfig is a specialized mock for testing the config generation
type MockSSOClientForConfig struct {
	mock.Mock
}

func (m *MockSSOClientForConfig) ListAccounts(accounts []types.AccountInfo, err error) {
	m.On("ListAccounts", mock.Anything, mock.Anything).Return(&sso.ListAccountsOutput{
		AccountList: accounts,
	}, err)
}

func (m *MockSSOClientForConfig) HasMorePages() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSSOClientForConfig) NextPage(ctx interface{}) (*sso.ListAccountsOutput, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sso.ListAccountsOutput), args.Error(1)
}

// TestShowFileDiff tests the showFileDiff function
func TestShowFileDiff(t *testing.T) {
	// Create two test files
	tempDir := t.TempDir()
	file1Path := tempDir + "/file1"
	file2Path := tempDir + "/file2"

	// Create test files
	err := os.WriteFile(file1Path, []byte("test content 1"), 0600)
	require.NoError(t, err)
	err = os.WriteFile(file2Path, []byte("test content 2"), 0600)
	require.NoError(t, err)

	// Call the function (it just prints to stdout, so we're mostly testing it doesn't crash)
	showFileDiff(file1Path, file2Path)

	// Test with nonexistent file
	showFileDiff(file1Path, tempDir+"/nonexistent")
	showFileDiff(tempDir+"/nonexistent", file2Path)
}

// TestMockConfigGeneration tests the mock configuration
func TestMockConfigGeneration(t *testing.T) {
	t.Skip("This test is fragile due to mocking limitations with AWS SDK pagination")
}

// MockConfigGenerator implements ConfigGeneratorIface for testing
type MockConfigGenerator struct {
	mock.Mock
}

func (m *MockConfigGenerator) ListAccountsWithClient(ssoClient SSOClient, token *string) ([]types.AccountInfo, error) {
	args := m.Called(ssoClient, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.AccountInfo), args.Error(1)
}

func (m *MockConfigGenerator) GetAccountRolesWithClient(ssoClient SSOClient, token *string, accountID string) ([]types.RoleInfo, error) {
	args := m.Called(ssoClient, token, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.RoleInfo), args.Error(1)
}

func (m *MockConfigGenerator) WriteSectionToConfig(configParser *configparser.ConfigParser, sectionName string, values map[string]string) error {
	args := m.Called(configParser, sectionName, values)
	return args.Error(0)
}

func (m *MockConfigGenerator) GenerateConfigFile(ssoClient SSOClient, token *string, configFile string, showDiff bool, appCfg *appconfig.Config) error {
	args := m.Called(ssoClient, token, configFile, showDiff, appCfg)
	return args.Error(0)
}

// TestGenerateAwsConfigFileWithMocks tests the AWS config file generation using mocks
func TestGenerateAwsConfigFileWithMocks(t *testing.T) {
	t.Skip("This test is skipped due to complexity in properly mocking AWS SSO pagination")
}
