package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	log.Info().Str("module", "config").Msgf("config init success")
	return &Config{}
}

func (c *Config) GetJiraConfig() *JiraConfig {

	return &JiraConfig{
		ReviewStatus: viper.GetString("PM_TOOL_JIRA_REVIEW_STATUS"),
		ExecStatus:   viper.GetString("PM_TOOL_JIRA_EXEC_STATUS"),
	}
}

// Config structs
type JiraConfig struct {
	ReviewStatus string
	ExecStatus   string
}
