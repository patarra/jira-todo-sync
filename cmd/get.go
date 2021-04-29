/*
Copyright © 2021 Jose Manuel Felguera Rodriguez <patarra@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/patarra/jira-todo-sync/jira"
	"github.com/patarra/jira-todo-sync/utils"
	"github.com/spf13/cobra"
	jira2 "gopkg.in/andygrunwald/go-jira.v1"
	"os"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets your assigned tasks from JIRA",
	Long:  `Gets your assigned tasks from JIRA`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := jira.GetJiraClient()
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
		issues, err := client.GetIssuesAssignedNotClosed()
		if len(issues) <= 0{
			// jira could return [] if the credentials are not valid
			utils.PrintInfoF("No results from JIRA.")
			utils.PrintInfoF("Jira returns and empty search result if the credentials are not valid, please check them")
		}
		printIssues(issues)
	},
}

func printIssues(issues []jira2.Issue){
	for _,issue := range issues {
		fmt.Printf("Issue: %s \n", issue.Key)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
