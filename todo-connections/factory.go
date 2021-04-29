package todo_connections

import (
	"fmt"
	"github.com/patarra/jira-todo-sync/config"
	"gopkg.in/andygrunwald/go-jira.v1"
)

type TodoManager interface {
	SyncIssues(issues jira.Issue) error
}

func GetTodoManager(name string) (TodoManager,error) {
	if name == "todoist"{
		cfg, _:=config.GetConfig()
		return newTodoistManager(cfg.Todoist)
	}
	return nil, fmt.Errorf("wrong todo manager id: %s",name)
}