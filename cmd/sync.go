/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/patarra/jira-todo-sync/jira"
	"github.com/patarra/jira-todo-sync/spi"
	"github.com/patarra/jira-todo-sync/utils"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tasks to your favourite todo app",
	Long:  "Sync tasks to your favourite todo app",
	Run: func(cmd *cobra.Command, args []string) {
		syncApp, _ := cmd.Flags().GetString("to")
		if !spi.IsTodoAppAvailable(syncApp) {
			utils.PrintErrorF("%s app is invalid", syncApp)
			cmd.Help()
		}
		// get issues from jira
		client, err := jira.GetJiraClient()
		if err != nil {
			utils.PrintError(err)
			os.Exit(1)
		}
		issues, err := client.GetIssuesAssignedNotClosed()
		if err != nil {
			utils.PrintError(err)
			os.Exit(1)
		}
		app, err := spi.GetTodoManager(syncApp)
		if err != nil {
			utils.PrintError(err)
			os.Exit(1)
		}
		err = app.SyncIssues(issues, spi.GetFlagValues(cmd, app.GetFlagsDescriptions()))
		if err != nil{
			utils.PrintError(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringP("to", "t", "", "Todo app -> (todoist)")
	_ = syncCmd.MarkFlagRequired("to")

	// set up flags for providers
	for _, i := range spi.GetAllTodoManagers() {
		flags := i.GetFlagsDescriptions()
		for _, f := range flags {
			syncCmd.Flags().StringP(f.Name, f.Shorthand, f.DefaultValue, f.Description)
			if f.Required {
				//ignore errors
				_ = syncCmd.MarkFlagRequired(f.Name)
			}
		}
	}
}
