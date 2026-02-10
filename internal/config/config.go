package config

import (
	"fmt"
	"os"
)

type Config struct {
	APIToken string
	TeamID   string
}

func Load() (*Config, error) {
	token := os.Getenv("CLICKUP_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("CLICKUP_API_TOKEN environment variable is required")
	}

	teamID := os.Getenv("CLICKUP_TEAM_ID")
	if teamID == "" {
		return nil, fmt.Errorf("CLICKUP_TEAM_ID environment variable is required")
	}

	return &Config{
		APIToken: token,
		TeamID:   teamID,
	}, nil
}
