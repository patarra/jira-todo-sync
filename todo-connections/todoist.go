package todo_connections

import (
	"fmt"
	todoistlib "github.com/kobtea/go-todoist/todoist"
	"github.com/patarra/jira-todo-sync/config"
	"gopkg.in/andygrunwald/go-jira.v1"
)

type todoist struct {
	client *todoistlib.Client
}

func newTodoistManager(config config.TodoistConfig) (TodoManager,error){
	//build client for todoist
	if !validateTodoistConfig(config){
		return nil, fmt.Errorf("todoist token must be present in the config")
	}

	client,err:=todoistlib.NewClient("", config.Token,"*","",nil)
	if err != nil {
		return nil,err
	}
	return &todoist{client: client}, nil
}

func validateTodoistConfig(config config.TodoistConfig) bool{
	return len(config.Token) > 0
}
func (t todoist) SyncIssues(issues jira.Issue) error{
	return nil
}