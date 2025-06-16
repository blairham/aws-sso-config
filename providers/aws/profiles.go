package aws

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bigkevmcd/go-configparser"
)

const AwsProfile = "AWS_PROFILE"

var Logger *log.Logger = log.New(os.Stderr, "", 0)

func validateAccountID(accountID, rootDir string) error {
	terragruntFile := filepath.Join(rootDir, "terragrunt.hcl")

	readFile, err := os.Open(terragruntFile)
	if err != nil {
		return fmt.Errorf("could not find terragrunt.hcl at root of git repo %s", rootDir)
	}

	defer func() {
		if closeErr := readFile.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	accountIDFromTerragrunt := ""
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineFields := strings.Fields(line)
		if len(lineFields) < 3 {
			continue
		}
		if lineFields[0] != "account_id" {
			continue
		}
		accountIDFromTerragrunt = strings.Trim(lineFields[2], "\"")
		break
	}
	if accountIDFromTerragrunt == "" {
		return fmt.Errorf("could not determine account id from %s", terragruntFile)
	}
	if accountIDFromTerragrunt != accountID {
		return fmt.Errorf("account id %s determined from profile did not match entry in terragrunt file %s", accountID, terragruntFile)
	}

	return nil
}

func getProfileFromRepoName(repo string) string {
	// For most repositories, the profile name matches the repo name
	// This function can be extended to handle special cases if needed
	return repo
}

// Exit if profile does not appear to be valid
func validateProfile(awsProfile, rootDir string) error {
	configFile, err := ConfigFile()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	awsConfig, err := configparser.NewConfigParserFromFile(configFile)
	if err != nil {
		return err
	}
	section := fmt.Sprintf("profile %s", awsProfile)

	if !awsConfig.HasSection(section) {
		return fmt.Errorf("could not find profile for %s", awsProfile)
	}
	accountID, err := awsConfig.Get(section, "sso_account_id")
	if err != nil {
		return fmt.Errorf("error parsing aws config %s: %s", awsProfile, err)
	}
	err = validateAccountID(accountID, rootDir)
	if err != nil {
		return err
	}

	Logger.Printf("Using profile %s (%s)", awsProfile, accountID)
	return nil
}

func GetProfile() (string, error) {
	val, present := os.LookupEnv(AwsProfile)
	if present {
		Logger.Printf("%s is already set to %s (potentially by direnv?), skipping setup", AwsProfile, val)
		return "", nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		_, err = os.Stat(filepath.Join(cwd, ".git"))
		if err == nil {
			break
		}
		cwd = filepath.Dir(cwd)
		if cwd == "/" {
			// If no git repository is found, return the default profile
			return "default", nil
		}
	}
	repoName := filepath.Base(cwd)
	profile := getProfileFromRepoName(repoName)
	err = validateProfile(profile, cwd)
	if err != nil {
		// If there's an error validating the profile (e.g., no remotes), return default
		return "default", nil
	}
	return profile, nil
}
