package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigDir is the directory where the config file is stored
var (
	BaseURL   = "https://dev.api.mockthis.io/api/v1"
	TokenFile = ".credentials"
	ConfigDir = ".mockthis"
)

// Data is the structure of the config file
type Data struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

// SaveConfig saves the config to the config file
func SaveConfig(filename, data string) error {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configPath, filename), []byte(data), 0600)
}

// LoadConfig loads the config from the config file
func LoadConfig(filename string) (*Data, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir, filename)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Data
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
