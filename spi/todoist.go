package spi

import (
	"context"
	"fmt"
	"github.com/patarra/jira-todo-sync/utils"
	"strings"

	todoistlib "github.com/kobtea/go-todoist/todoist"
	"github.com/patarra/jira-todo-sync/config"
	"gopkg.in/andygrunwald/go-jira.v1"
)

type todoist struct {
	client *todoistlib.Client
}


func init() {
	err := RegisterManager("todoist", &todoist{})
	if err != nil {
		panic(err)
	}
}

func (t *todoist) GetFlagsDescriptions() []FlagDescriptor {
	var result []FlagDescriptor
	result = append(result, FlagDescriptor{
		Name:         "todoist-label",
		Shorthand:    "",
		DefaultValue: "jira",
		Description:  "Only items with this label will be considered. Also, all items from jira will be created with this flag attached",
		Required:     false,
	})
	result = append(result, FlagDescriptor{
		Name:         "todoist-project",
		Shorthand:    "",
		DefaultValue: "Inbox",
		Description:  "Items synced from jira will be created in this project",
		Required:     false,
	})
	return result
}

func (t *todoist) SyncIssues(issues []jira.Issue, flagValues map[string]string ) error {
	_, err := t.getTodoistClient()
	if err != nil{
		return err
	}
	existentItems, err := t.filterByLabel(flagValues["todoist-label"])
	if err!=nil{
		return err
	}
	if err != nil {
		return err
	}
	for _, i := range issues {
		if existJiraTicketInTodoist(existentItems, i.Key) {
			utils.PrintInfoF("%s already exists", i.Key)
		} else {
			err = t.createTodoistTask(i, flagValues["todoist-project"], flagValues["todoist-label"])
			if err != nil {
				utils.PrintErrorF("%s failed to create in todoist: %s", i.Key, err)
			}else{
				utils.PrintSuccessF("%s created", i.Key)
			}
		}
	}
	return t.client.Commit(context.Background())
}

func (t *todoist) createTodoistTask(i jira.Issue, projectName string, labelName string ) error{
	project := t.findProject(projectName)
	label, _ := t.findLabel(labelName)
	opts := todoistlib.NewItemOpts{
		Labels: []todoistlib.ID{label.ID},
	}
	if project != nil {
		opts.ProjectID=project.ID
	}
	item,_ := todoistlib.NewItem(fmt.Sprintf("%s - %s", i.Key, i.Fields.Summary), &opts)
	_, err := t.client.Item.Add(*item)
	return err
}

func (t *todoist) getTodoistClient() (*todoistlib.Client, error) {
	if t.client == nil {
		cfg, err := config.GetConfig()
		if len(cfg.Todoist.Token) <= 0 {
			return nil, fmt.Errorf("todoist token must be present in the config")
		}

		client, err := todoistlib.NewClient("", cfg.Todoist.Token, "*", "", nil)
		if err != nil {
			return nil, err
		}
		ctx := context.Background()
		err = client.FullSync(ctx, []todoistlib.Command{})
		if err != nil {
			return nil, err
		}
		t.client = client
	}
	return t.client, nil
}

func existJiraTicketInTodoist(items []todoistlib.Item, identifier string) bool {
	for _, i := range items {
		if strings.HasPrefix(i.Content, identifier) {
			return true
		}
	}
	return false
}

func (t *todoist) filterByLabel(labelName string) ([]todoistlib.Item, error) {
	label, err := t.findLabel(labelName)
	if err != nil {
		return nil, err
	}
	var result []todoistlib.Item
	for _, i := range t.client.Item.GetAll() {
		for _, j := range i.Labels {
			if j == label.ID {
				result = append(result, i)
				break
			}
		}
	}
	return result, nil
}

func (t *todoist) findLabel(labelName string) (*todoistlib.Label, error) {
	labels := t.client.Label.GetAll()
	for _, l := range labels {
		if l.Name == labelName {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("todoist: label %s not found", labelName)
}

func (t *todoist) findProject(projectName string) *todoistlib.Project {
	project := t.client.Project.FindByName(projectName)
	if len(project) > 0 {
		return &project[0]
	}
	return nil
}