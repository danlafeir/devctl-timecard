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

	"github.com/danlafeir/devctl-timecard/cmd/timecard"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version:      "1.0",
	Use:          "timecard",
	Short:        "commands to manage your timecard",
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var completionCmd = &cobra.Command{
	Use:    "completion [bash|zsh|fish|powershell]",
	Short:  "Generate shell completion scripts",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Println("Please specify a shell: bash, zsh, fish, or powershell")
			os.Exit(1)
		}
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			cmd.Println("Unsupported shell type.")
			os.Exit(1)
		}
	},
}

func init() {
	// Hide the help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	// Disable the completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Add commands
	rootCmd.AddCommand(timecard.AddEntryCmd())
	rootCmd.AddCommand(timecard.ConfigureCmd())
	rootCmd.AddCommand(timecard.GetWeekCmd())

	// Hide completion command if it was already registered
	if compCmd, _, _ := rootCmd.Find([]string{"completion"}); compCmd != nil {
		compCmd.Hidden = true
	}
}
