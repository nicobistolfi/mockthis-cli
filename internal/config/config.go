package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	BaseURL   = "https://dev.api.mockthis.io/api/v1"
	TokenFile = ".credentials"
	ConfigDir = ".mockthis"
)

type ConfigData struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func SaveConfig(filename, data string) error {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configPath, filename), []byte(data), 0600)
}

func LoadConfig(filename string) (*ConfigData, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir, filename)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ConfigData
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
