package spi

import (
	"fmt"
	"github.com/spf13/cobra"

	"gopkg.in/andygrunwald/go-jira.v1"
)

type FlagDescriptor struct {
	Name         string
	Shorthand    string
	DefaultValue string
	Description  string
	Required     bool
}

type TodoManager interface {
	GetFlagsDescriptions() []FlagDescriptor
	// todo: create a common jira object for issues
	SyncIssues(issues []jira.Issue, flagValues map[string]string) error
}

var providers = make(map[string]*TodoManager)

func RegisterManager(name string, manager TodoManager) error {
	if providers[name] != nil {
		return fmt.Errorf("provider %s already registered", name)
	}
	providers[name] = &manager
	return nil
}

func GetTodoManager(name string) (TodoManager, error) {

	provider := providers[name]
	if provider != nil {
		return *provider, nil
	}
	return nil, fmt.Errorf("wrong todo manager id: %s", name)
}

func GetFlagValues(cmd *cobra.Command, descriptors []FlagDescriptor) map[string]string {
	result := make(map[string]string)
	for _, f := range descriptors {
		value, err := cmd.Flags().GetString(f.Name)
		// if there is a value, add this value
		if err == nil && value != "" {
			result[f.Name] = value
		}else{
			result[f.Name] = f.DefaultValue
		}
	}
	return result
}

func GetAllTodoManagers() []TodoManager {
	var result []TodoManager
	for _, i := range providers {
		// check usage of * vs &
		result = append(result, *i)
	}
	return result
}

func IsTodoAppAvailable(name string) bool {
	provider := providers[name]
	if provider != nil {
		return true
	}
	return false
}
