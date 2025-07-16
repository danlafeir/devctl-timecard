package tempo

import (
	"fmt"
	"os"

	"github.com/danlafeir/devctl-tempo/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ConfigureCmd() *cobra.Command {
	var apiToken string
	var accountId string

	configureCmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure your Tempo API token and default issue ID",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig()
			viper.ReadInConfig()

			configureApiToken(apiToken)
			configureAccountId(accountId)
			// Get accountId from viper after it's been set
			configuredAccountId := viper.GetString("tempo." + ACCOUNT_ID)
			if configuredAccountId == "" {
				configuredAccountId = accountId
			}
			configureIssueId(configuredAccountId)

			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Failed to save config:", err)
				os.Exit(1)
			}
			fmt.Println("Configuration saved successfully.")
		},
	}
	configureCmd.Flags().StringVar(&apiToken, "token", "", "Tempo API token")
	configureCmd.Flags().StringVar(&accountId, "account-id", "", "Tempo account ID")
	return configureCmd
}

func TimesheetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "timesheet",
		Short:   "Submit a Tempo timesheet",
		Example: "devctl tempo timesheet",
		RunE: func(cmd *cobra.Command, args []string) error {
			bearerToken := fetchBearerToken()
			accountId, issueId := fetchConfig()
			startOfWeek := requestDayOfWeek()
			devTime, ptoTime, meetingTime := requestTimeInput()

			if err := api.SendWorklog(api.DevWorkType, devTime, startOfWeek, bearerToken, accountId, issueId); err != nil {
				return err
			}
			if err := api.SendWorklog(api.PtoWorkType, ptoTime, startOfWeek, bearerToken, accountId, issueId); err != nil {
				return err
			}
			if err := api.SendWorklog(api.MeetingWorkType, meetingTime, startOfWeek, bearerToken, accountId, issueId); err != nil {
				return err
			}

			fmt.Println("✅ All time entries submitted successfully!")
			return nil
		},
	}
	return cmd
}
