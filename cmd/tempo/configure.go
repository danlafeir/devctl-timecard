package tempo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/spf13/viper"
)

const ACCOUNT = "tempo"
const API_TOKEN_LOCATION = "cli.dpctl.tempo"
const ACCOUNT_ID = "tempoAccountId"
const ISSUE_ID = "tempoIssueId"
var configPath string 

func configureApiToken(apiToken string) string {
	token := strings.TrimSpace(apiToken)
	if token == "" {
		fmt.Print("Enter your Tempo API token: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		token = strings.TrimSpace(scanner.Text())
	}

	if token == "" {
		fmt.Println("Token cannot be empty. Re-run the configure command")
		os.Exit(1)
	}
	
	if err := secrets.DefaultSecrets.Write(API_TOKEN_LOCATION, token); err != nil {
		fmt.Println("Failed to write token to keychain:", err)
		os.Exit(1)
	}
	fmt.Println("Tempo API token saved securely to keychain.")
	return token
}

func configureAccountId(accountId string) {
	if accountId == "" {
		fmt.Printf("We need your Tempo Account Id. This is associated with your profile in JIRA.\n")
		fmt.Printf("Do a 'people search' to find your id.\n")
		fmt.Printf("Add Tempo Account Id here: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		accountId = strings.TrimSpace(scanner.Text())
		fmt.Print("\n")
	}
	if accountId == "" {
		fmt.Println("Account ID cannot be empty.")
		os.Exit(1)
	}
	viper.Set(ACCOUNT_ID, accountId)
}

func configureIssueId(issueId string) {
	id := strings.TrimSpace(issueId)
	if id == "" {
		fmt.Print("Enter your default Issue ID: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		id = strings.TrimSpace(scanner.Text())
		fmt.Print("\n")
	}
	if id == "" {
		fmt.Println("Issue ID cannot be empty.")
		os.Exit(1)
	}
	viper.Set(ISSUE_ID, id)

}

func getConfigPath() string {
	if configPath != "" {
		return configPath
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}
	return filepath.Join(homeDir, ".devctl", "plugins", "config-tempo.json")
}

func initConfig() {
	configFilePath := getConfigPath()
	configDir := filepath.Dir(configFilePath)
	configName := strings.TrimSuffix(filepath.Base(configFilePath), filepath.Ext(configFilePath))
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Failed to create config directory:", err)
	}
	
	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)
}

func fetchConfig() (accountId string, issueId string) {
	initConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if !viper.IsSet(ACCOUNT_ID) {
		configureAccountId("")
	} else {
		accountId = viper.GetString(ACCOUNT_ID)
	}

	if !viper.IsSet(ISSUE_ID) {
		configureIssueId("")
	} else {
		issueId = viper.GetString(ISSUE_ID)
	}

	return
}

func fetchBearerToken() string {
	bearerToken, err := secrets.DefaultSecrets.Read(API_TOKEN_LOCATION)

	if bearerToken == "" || err != nil  {
		return configureApiToken(bearerToken)
	}
	return bearerToken
}