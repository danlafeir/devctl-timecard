package tempo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ACCOUNT = "tempo"
const ACCESS_GROUP = "cli.dpctl.tempo"
const ACCOUNT_ID = "tempoAccountId"
const ISSUE_ID = "tempoIssueId"

func getConfigPath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}
	return filepath.Join(homeDir, ".devctl", "plugins", "config-tempo.json")
}

func initConfig(customPath string) {
	configFilePath := getConfigPath(customPath)
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

type Config struct {
	TempoAPIToken string `json:"tempo_api_token"`
}

var apiToken string
var issueId string
var accountId string
var configPath string

func ConfigureCmd() *cobra.Command {
	configureCmd := &cobra.Command{
	Use:   "configure",
	Short: "Configure your Tempo API token and default issue ID",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize config
		initConfig(configPath)
		viper.ReadInConfig()

		// API Token logic
		token := strings.TrimSpace(apiToken)
		if token == "" {
			fmt.Print("Enter your Tempo API token: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			token = strings.TrimSpace(scanner.Text())
		}
		if token == "" {
			fmt.Println("Token cannot be empty.")
			os.Exit(1)
		}
		
		if err := secrets.DefaultSecrets.Write("tempo_api_token", []byte(token)); err != nil {
			fmt.Println("Failed to write token to keychain:", err)
			os.Exit(1)
		}
		fmt.Println("Tempo API token saved securely to keychain.")

		// Account ID logic
		accId := strings.TrimSpace(accountId)
		if accId == "" {
			fmt.Print("Enter your Tempo Account ID: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			accId = strings.TrimSpace(scanner.Text())
		}
		if accId == "" {
			fmt.Println("Account ID cannot be empty.")
			os.Exit(1)
		}
		viper.Set(ACCOUNT_ID, accId)

		// Issue ID logic
		id := strings.TrimSpace(issueId)
		if id == "" {
			fmt.Print("Enter your default Issue ID: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			id = strings.TrimSpace(scanner.Text())
		}
		if id == "" {
			fmt.Println("Issue ID cannot be empty.")
			os.Exit(1)
		}
		viper.Set(ISSUE_ID, id)
		
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Failed to save config:", err)
			os.Exit(1)
		}
		fmt.Println("Configuration saved successfully.")
	},
	}
	configureCmd.Flags().StringVar(&apiToken, "token", "", "Tempo API token")
	configureCmd.Flags().StringVar(&issueId, "issue-id", "", "Default issue ID")
	configureCmd.Flags().StringVar(&accountId, "account-id", "", "Tempo account ID")
	configureCmd.Flags().StringVar(&configPath, "config-path", "", "Path to config file (default: $HOME/.devctl/plugins/config-tempo.json)")
	return configureCmd
}

func FetchConfig() (accountId string, issueId string) {
	initConfig("") // Use default path
	
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if !viper.IsSet(ACCOUNT_ID) {
		fmt.Printf("We need your Tempo Account Id. This is associated with your profile in JIRA.\n")
		fmt.Printf("Do a 'people search' to find your id.\n")
		fmt.Printf("Add Tempo Account Id here: ")

		if _, err := fmt.Scan(&accountId); err != nil {
			log.Fatal(err)
		}
		viper.Set(ACCOUNT_ID, accountId)
		viper.WriteConfig()

		fmt.Printf("\n")
	} else {
		accountId = viper.GetString(ACCOUNT_ID)
	}

	if !viper.IsSet(ISSUE_ID) {
		fmt.Printf("Add your Issue Id here: ")

		if _, err := fmt.Scan(&issueId); err != nil {
			log.Fatal(err)
		}
		viper.Set(ISSUE_ID, issueId)
		viper.WriteConfig()

		fmt.Printf("\n")
	} else {
		issueId = viper.GetString(ISSUE_ID)
	}

	return
}

func FetchBearerToken() string {
	bearerToken := checkForApiToken()
	if bearerToken == "" {
		return addApiTokenToKeychain()
	}
	return bearerToken
}

func checkForApiToken() string {
	token, err := secrets.DefaultSecrets.Read("tempo_api_token")
	if err != nil {
		return ""
	}
	return string(token)
}

func addApiTokenToKeychain() string {
	fmt.Printf("Grab an API token from https://grainger.atlassian.net/plugins/servlet/ac/io.tempo.jira/tempo-app#!/configuration/api-integration \n")
	fmt.Printf("Add API Token here: ")

	var apiToken string
	if _, err := fmt.Scan(&apiToken); err != nil {
		log.Fatal(err)
	}

	if err := secrets.DefaultSecrets.Write("tempo_api_token", []byte(apiToken)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n")
	return apiToken
}
