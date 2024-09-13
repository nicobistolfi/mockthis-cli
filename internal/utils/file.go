package utils

import (
	"os"
	"path/filepath"
)

func LoadFile(path string) (string, error) {
	// Check if file exists
	exists, err := FileExists(path)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", os.ErrNotExist
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteFile(path string, content string) error {
	// make directory if not exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
