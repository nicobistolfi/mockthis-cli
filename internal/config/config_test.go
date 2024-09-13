package config

import (
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
	testConfig := &Data{
		Token: "test_token",
		Email: "test@example.com",
	}

	// Test SaveConfig
	err = SaveConfig(TokenFile, testConfig)
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
