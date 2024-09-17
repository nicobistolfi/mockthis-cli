package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

// UpdateEndpointCmd is the command to update an existing mock endpoint
var UpdateEndpointCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update an existing mock endpoint",
	Args:  cobra.ExactArgs(1),
	Run:   updateEndpoint,
}

func updateEndpoint(cmd *cobra.Command, args []string) {
	configData, err := config.LoadConfig(config.TokenFile)
	if err != nil {
		fmt.Println("You need to login first.")
		return
	}
	token := configData.Token

	id := args[0]
	status := promptForInput("Enter new HTTP status code: ")
	responseContentType := promptForInput("Enter new response Content-Type: ")
	charset := promptForInput("Enter new charset: ")
	responseBody := promptForInput("Enter new response body: ")

	endpointData := map[string]interface{}{
		"status":              status,
		"responseContentType": responseContentType,
		"charset":             charset,
		"responseBody":        responseBody,
	}
	jsonData, _ := json.Marshal(endpointData)

	req, _ := http.NewRequest("PATCH", config.BaseURL+"/endpoints/"+id, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error updating endpoint:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to update endpoint. Status:", resp.Status)
		return
	}

	fmt.Println("Endpoint updated successfully!")
}
