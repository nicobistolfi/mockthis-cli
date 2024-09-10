package main

import (
	"testing"

	"github.com/nicobistolfi/mockthis-cli/internal/commands"
	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	if rootCmd.Use != "mockthis" {
		t.Errorf("Expected root command Use to be 'mockthis', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "MockThis - A CLI for managing mock API endpoints" {
		t.Errorf("Expected root command Short description to be 'MockThis - A CLI for managing mock API endpoints', got '%s'", rootCmd.Short)
	}
}

func TestCommandInitialization(t *testing.T) {
	expectedCommands := map[string]*cobra.Command{
		"login":    commands.LoginCmd,
		"register": commands.RegisterCmd,
		"create":   commands.CreateEndpointCmd,
		"list":     commands.ListEndpointsCmd,
		"get":      commands.GetEndpointCmd,
		"update":   commands.UpdateEndpointCmd,
		"delete":   commands.DeleteEndpointCmd,
	}

	for name, expectedCmd := range expectedCommands {
		t.Run(name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{name})
			if err != nil {
				t.Errorf("Expected to find command '%s', but got error: %v", name, err)
			}
			if cmd != expectedCmd {
				t.Errorf("Expected command '%s' to be %v, but got %v", name, expectedCmd, cmd)
			}
		})
	}
}

func TestMainFunction(t *testing.T) {
	// Since the main function calls rootCmd.Execute(), which is difficult to test directly,
	// we'll just ensure that rootCmd is properly initialized.
	if rootCmd == nil {
		t.Error("Expected rootCmd to be initialized, but it's nil")
	}

	if len(rootCmd.Commands()) != 7 {
		t.Errorf("Expected rootCmd to have 7 subcommands, but got %d", len(rootCmd.Commands()))
	}
}
