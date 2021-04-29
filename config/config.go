package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type JiraConfig struct {
	Server   string `mapstructure:"server"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type TodoistConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Jira    JiraConfig    `mapstructure:"jira"`
	Todoist TodoistConfig `mapstructure:"todoist"`
}

var instance Config
var once sync.Once
var initialised = false

func InitConfig(cfgFile string) (*Config, error) {
	var onceErr error = nil
	once.Do(func() {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			onceErr = errors.New(fmt.Sprintf("couldn't load config from %s: %s\n", cfgFile, err))
		}
		if err := viper.Unmarshal(&instance); err != nil {
			onceErr = errors.New(fmt.Sprintf("couldn't read config: %s\n", err))
		}
		initialised = true
	})
	if onceErr == nil {
		return &instance, nil
	} else {
		return nil, onceErr
	}
}

func GetConfig() (*Config, error) {
	if !initialised {
		return nil, errors.New("config is not initialised yet, please call InitConfig(cfgFile)")
	}
	return &instance, nil
}

