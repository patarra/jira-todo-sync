package spi

import (
	"fmt"
	"gopkg.in/andygrunwald/go-jira.v1"
)

type FlagDescriptor struct{
	Name string
	Shorthand string
	DefaultValue string
	Description string
	Required bool
}

type TodoManager interface {
	GetFlagsDescriptions() []FlagDescriptor
	// todo: create a common jira object for issues
	SyncIssues(issues []jira.Issue, flags []FlagDescriptor) error
}

var providers = make(map[string]*TodoManager)

func RegisterManager(name string, manager TodoManager) error{
	if providers[name] != nil{
		return fmt.Errorf("provider %s already registered", name)
	}
	providers[name] = &manager
	return nil
}

func GetTodoManager(name string) (*TodoManager, error) {

	provider := providers[name]
	if provider != nil{
		return provider,nil
	}
	return nil, fmt.Errorf("wrong todo manager id: %s", name)
}

func GetAllTodoManagers() []TodoManager {
	var result []TodoManager
	for _,i := range providers{
		// check usage of * vs &
		result = append(result,*i)
	}
	return result
}

