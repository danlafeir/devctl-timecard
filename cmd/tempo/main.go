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
	var issueId string
	var accountId string
	
	configureCmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure your Tempo API token and default issue ID",
		Run: func(cmd *cobra.Command, args []string) {
			viper.ReadInConfig()

			configureApiToken(apiToken)
			configureAccountId(accountId)
			configureIssueId(issueId)
			
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
	return configureCmd
}

func TimesheetCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:     "timesheet",
		Short:   "Submit a Tempo timesheet",
		Example: "dpctl tempo timesheet",
		RunE: func(cmd *cobra.Command, args []string) error {
			bearerToken := fetchBearerToken()
			accountId, issueId := fetchConfig()
			startOfWeek := requestDayOfWeek()
			devTime, ptoTime, supportTime, meetingTime := requestTimeInput()
			if dryRun {
				fmt.Println(devTime, ptoTime, supportTime, meetingTime)
				fmt.Println("Skipping posting time to Tempo")
			} else {
				api.SendWorklog(api.DevWorkType, devTime, startOfWeek, bearerToken, accountId, issueId)
				api.SendWorklog(api.PtoWorkType, ptoTime, startOfWeek, bearerToken, accountId, issueId)
				api.SendWorklog(api.SupportWorkType, supportTime, startOfWeek, bearerToken, accountId, issueId)
				api.SendWorklog(api.MeetingWorkType, meetingTime, startOfWeek, bearerToken, accountId, issueId)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't post time to tempo and print times instead")
	return cmd
}
