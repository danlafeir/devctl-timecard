/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/danlafeir/devctl-tempo/cmd/tempo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version:      "1.0",
	Use:          "tempo",
	Short:        "Tempo timesheet management commands",
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(tempo.TimesheetCmd())
	rootCmd.AddCommand(tempo.ConfigureCmd())
	rootCmd.AddCommand(tempo.GetWeekCmd())
}
