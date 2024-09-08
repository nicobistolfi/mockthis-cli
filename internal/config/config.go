package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	BaseURL   = "https://dev.api.mockthis.io/api/v1"
	TokenFile = ".token"
	EmailFile = ".email"
	ConfigDir = ".mockthis"
)

func SaveConfig(filename, data string) error {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(configPath, filename), []byte(data), 0600)
}

func LoadConfig(filename string) (string, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigDir, filename)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
