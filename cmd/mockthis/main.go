package main

import (
	"fmt"
	"os"

	"github.com/nicobistolfi/mockthis-cli/internal/commands"
	"github.com/spf13/cobra"
)

var (
	version     = "dev"
	commit      = "none"
	date        = "unknown"
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "mockthis",
	Short: "MockThis - A CLI for managing mock API endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Date: %s\n", date)
			return
		} else {
			_ = cmd.Help()
		}
	},
}

func init() {

	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.RegisterCmd)
	rootCmd.AddCommand(commands.CreateEndpointCmd)
	rootCmd.AddCommand(commands.ListEndpointsCmd)
	rootCmd.AddCommand(commands.GetEndpointCmd)
	rootCmd.AddCommand(commands.UpdateEndpointCmd)
	rootCmd.AddCommand(commands.DeleteEndpointCmd)

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
