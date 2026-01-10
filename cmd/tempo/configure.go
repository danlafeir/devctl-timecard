package tempo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danlafeir/devctl-timecard/api"
	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/spf13/viper"
)

const ACCOUNT = "tempo"
const API_TOKEN_LOCATION = "cli.devctl.tempo"
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
		fmt.Print("Add Tempo Account Id here: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		accountId = strings.TrimSpace(scanner.Text())
		fmt.Print("\n")
	}
	if accountId == "" {
		fmt.Println("Account ID cannot be empty.")
		os.Exit(1)
	}
	viper.Set("tempo."+ACCOUNT_ID, accountId)
}

func configureIssueId(accountId string) {
	if accountId == "" {
		fmt.Println("Account ID is required to fetch recent issue ID. Please configure account ID first.")
		os.Exit(1)
	}

	bearerToken := fetchBearerToken()
	fmt.Print("Fetching recent issue ID from Tempo API...\n")
	recentIssueId, err := api.GetRecentIssueId(accountId, bearerToken)
	if err != nil {
		fmt.Printf("Failed to fetch recent issue ID: %v\n", err)
		fmt.Print("Enter your default Issue ID manually: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		id := strings.TrimSpace(scanner.Text())
		if id == "" {
			fmt.Println("Issue ID cannot be empty.")
			os.Exit(1)
		}
		viper.Set("tempo."+ISSUE_ID, id)
		return
	}

	id := strconv.Itoa(recentIssueId)
	fmt.Printf("Found recent issue ID: %s\n", id)
	viper.Set("tempo."+ISSUE_ID, id)
}

func getConfigPath() string {
	if configPath != "" {
		return configPath
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}

	// Detect binary name to determine config location
	// If running as "timecard" (standalone), use .timecard/config.yaml
	// Otherwise (devctl-timecard), use .devctl/config.yaml
	var execName string
	if len(os.Args) > 0 {
		execPath := os.Args[0]
		execName = filepath.Base(execPath)
		// Remove any extensions and check if it's the standalone binary
		execName = strings.TrimSuffix(execName, filepath.Ext(execName))
	}

	if execName == "timecard" {
		timecardConfigDir := filepath.Join(homeDir, ".timecard")
		// Ensure .timecard directory exists for standalone binary
		if err := os.MkdirAll(timecardConfigDir, 0755); err != nil {
			log.Fatal("Failed to create .timecard config directory:", err)
		}
		return filepath.Join(timecardConfigDir, "config.yaml")
	}

	// Default to devctl config location
	return filepath.Join(homeDir, ".devctl", "config.yaml")
}

func initConfig() {
	configFilePath := getConfigPath()
	configDir := filepath.Dir(configFilePath)
	configName := strings.TrimSuffix(filepath.Base(configFilePath), filepath.Ext(configFilePath))

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Failed to create config directory:", err)
	}

	// Create config file if it doesn't exist (with empty tempo structure)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Create empty config file with tempo key structure (tempo.* keys are used for Tempo API settings)
		emptyConfig := []byte("tempo:\n")
		if err := os.WriteFile(configFilePath, emptyConfig, 0644); err != nil {
			log.Fatal("Failed to create config file:", err)
		}
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
}

func fetchConfig() (accountId string, issueId string) {
	initConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if !viper.IsSet("tempo." + ACCOUNT_ID) {
		configureAccountId("")
	} else {
		accountId = viper.GetString("tempo." + ACCOUNT_ID)
	}

	if !viper.IsSet("tempo." + ISSUE_ID) {
		configureIssueId(accountId)
	} else {
		issueId = viper.GetString("tempo." + ISSUE_ID)
	}

	return
}

func fetchBearerToken() string {
	bearerToken, err := secrets.DefaultSecrets.Read(API_TOKEN_LOCATION)

	if bearerToken == "" || err != nil {
		return configureApiToken("")
	}
	return bearerToken
}
