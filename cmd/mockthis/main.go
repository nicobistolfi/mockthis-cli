package main

import (
	"fmt"
	"os"

	"github.com/nicobistolfi/mockthis-cli/internal/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mockthis",
	Short: "MockThis - A CLI for managing mock API endpoints",
}

func init() {
	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.RegisterCmd)
	rootCmd.AddCommand(commands.CreateEndpointCmd)
	rootCmd.AddCommand(commands.ListEndpointsCmd)
	rootCmd.AddCommand(commands.UpdateEndpointCmd)
	rootCmd.AddCommand(commands.DeleteEndpointCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
