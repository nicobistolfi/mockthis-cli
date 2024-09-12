package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadConfig(t *testing.T) {
	// Setup: Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the ConfigDir for testing
	originalConfigDir := ConfigDir
	ConfigDir = filepath.Base(tempDir)
	defer func() { ConfigDir = originalConfigDir }()

	// Override HOME environment variable
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test data
	testConfig := Data{
		Token: "test_token",
		Email: "test@example.com",
	}
	testData, _ := json.Marshal(testConfig)

	// Test SaveConfig
	err = SaveConfig(TokenFile, string(testData))
	if err != nil {
		t.Errorf("SaveConfig failed: %v", err)
	}

	// Test LoadConfig
	loadedConfig, err := LoadConfig(TokenFile)
	if err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	// Compare saved and loaded data
	if loadedConfig.Token != testConfig.Token || loadedConfig.Email != testConfig.Email {
		t.Errorf("Loaded config does not match saved config. Got %+v, want %+v", loadedConfig, testConfig)
	}
}

func TestSaveConfigError(t *testing.T) {
	// Test saving to an invalid directory
	err := SaveConfig("/invalid/path/file", "test data")
	if err == nil {
		t.Error("Expected an error when saving to an invalid directory, but got nil")
	}
}

func TestLoadConfigError(t *testing.T) {
	// Test loading a non-existent file
	_, err := LoadConfig("non_existent_file")
	if err == nil {
		t.Error("Expected an error when loading a non-existent file, but got nil")
	}

	// Test loading an invalid JSON
	tempDir, _ := os.MkdirTemp("", "config_test")
	defer os.RemoveAll(tempDir)

	originalConfigDir := ConfigDir
	ConfigDir = filepath.Base(tempDir)
	defer func() { ConfigDir = originalConfigDir }()

	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", os.Getenv("HOME"))

	invalidJSON := []byte(`{"token": "test", "email": }`) // Invalid JSON

	err = SaveConfig(TokenFile, string(invalidJSON))
	if err != nil {
		t.Error("Could not save invalid JSON")
	}

	_, err = LoadConfig(TokenFile)
	if err == nil {
		t.Error("Expected an error when loading invalid JSON, but got nil")
	}
}
