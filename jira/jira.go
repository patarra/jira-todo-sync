package jira

import (
	"errors"
	"fmt"
	"github.com/patarra/jira-todo-sync/config"
	"gopkg.in/andygrunwald/go-jira.v1"
	"sync"
)

type Helper struct {
	client *jira.Client
}

var instance *Helper
var lock = &sync.Mutex{}

func validateConfig(configJira *config.Config) bool {
	return len(configJira.Jira.Server) > 0
}

func GetJiraClient() (*Helper, error) {
	var jiraClient *jira.Client
	var err error
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		var cfg *config.Config
		cfg, err = config.GetConfig()
		if err == nil {
			if validateConfig(cfg) {
				if cfg.Jira.User != "" {
					tp := jira.BasicAuthTransport{
						Username: cfg.Jira.User,
						Password: cfg.Jira.Password,
					}
					// TODO: Verify error
					jiraClient, err = jira.NewClient(tp.Client(), cfg.Jira.Server)
				} else {
					jiraClient, err = jira.NewClient(nil, cfg.Jira.Server)
				}
				instance = &Helper{
					client: jiraClient,
				}
			} else {
				err = errors.New(fmt.Sprintf("Invalid configuration for jira, please ensure that you have specified a server"))
			}
		}
	}
	return instance, err
}

// Gets a list of issues assigned to the current user not closed
func (v Helper) GetIssuesAssignedNotClosed() ([]jira.Issue, error) {
	return v.searchIssues("assignee = currentUser() and status not in (\"Closed\")")
}

// Gets a list of all issues assigned to the user
func (v Helper) GetAllIssuesAssigned() ([]jira.Issue, error) {
	return v.searchIssues("assignee = currentUser()")
}

// will implement pagination of api and get all the issues.
// Jira API has limitation as to maxResults it can return at one time.
// You may have usecase where you need to get all the issues according to jql
// This is where this example comes in.
func (v Helper) searchIssues(searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
		}

		chunk, resp, err := v.client.Issue.Search(searchString, opt)
		if err != nil {
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}
		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= total {
			return issues, nil
		}
	}

}
