package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/nicobistolfi/mockthis-cli/internal/utils"
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
func SaveConfig(filename string, data *Data) error {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir, filename)
	jsonData, err := utils.ToYAML(data)
	if err != nil {
		return err
	}
	return utils.WriteFile(configPath, jsonData)
}

// LoadConfig loads the config from the config file
func LoadConfig(filename string) (*Data, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir, filename)
	data, err := utils.LoadFile(configPath)
	if err != nil {
		return nil, err
	}

	var parsedData map[string]interface{}

	if utils.IsJSON(data) {
		parsedDataJSON, err := utils.ParseJSON(data)
		parsedData = parsedDataJSON
		if err != nil {
			return nil, err
		}
	}

	if utils.IsYAML(data) {
		parsedDataYAML, err := utils.ParseYAML(data)
		parsedData = parsedDataYAML
		if err != nil {
			return nil, err
		}
	}

	if parsedData == nil {
		return nil, errors.New("config file is not a valid JSON or YAML")
	}

	if _, ok := parsedData["token"]; !ok {
		return nil, errors.New("token not found in config file")
	}

	if _, ok := parsedData["email"]; !ok {
		return nil, errors.New("email not found in config file")
	}

	config := &Data{
		Token: parsedData["token"].(string),
		Email: parsedData["email"].(string),
	}

	return config, nil
}
