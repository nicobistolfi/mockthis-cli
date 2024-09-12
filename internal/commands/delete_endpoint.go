package commands

import (
	"fmt"
	"net/http"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

// DeleteEndpointCmd is the command to delete an existing mock endpoint
var DeleteEndpointCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an existing mock endpoint",
	Args:  cobra.ExactArgs(1),
	Run:   deleteEndpoint,
}

func deleteEndpoint(cmd *cobra.Command, args []string) {
	configData, err := config.LoadConfig(config.TokenFile)
	if err != nil {
		fmt.Println("You need to login first.")
		return
	}
	token := configData.Token

	id := args[0]

	req, _ := http.NewRequest("DELETE", config.BaseURL+"/endpoints/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error deleting endpoint:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Failed to delete endpoint. Status:", resp.Status)
		return
	}

	fmt.Println("Endpoint deleted successfully!")
}
