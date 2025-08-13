package tempo

import (
	"fmt"
	"github.com/danlafeir/devctl-tempo/api"
	"github.com/spf13/cobra"
)


func TimesheetCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:     "timesheet",
		Short:   "Submit a Tempo timesheet",
		Example: "dpctl tempo timesheet",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tempoTime(dryRun)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't post time to tempo and print times instead")
	return cmd
}

func tempoTime(dryRun bool) error {
	bearerToken := FetchBearerToken()
	accountId, issueId := FetchConfig()
	startOfWeek := RequestDayOfWeek()
	devTime, ptoTime, supportTime, meetingTime := RequestTimeInput()
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
}