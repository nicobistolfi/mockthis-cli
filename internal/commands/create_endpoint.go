package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

var CreateEndpointCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new mock endpoint",
	Run:   createEndpoint,
}

func createEndpoint(cmd *cobra.Command, args []string) {
	configData, err := config.LoadConfig(config.TokenFile)
	if err != nil {
		fmt.Println("You need to login first.")
		return
	}
	token := configData.Token

	httpStatus := promptForInput("Enter HTTP status code: ")
	responseContentType := promptForInput("Enter response Content-Type: ")
	charset := promptForInput("Enter charset: ")
	responseBody := promptForInput("Enter response body: ")

	endpointData := map[string]interface{}{
		"httpStatus":          httpStatus,
		"responseContentType": responseContentType,
		"charset":             charset,
		"responseBody":        responseBody,
	}
	jsonData, _ := json.Marshal(endpointData)

	req, _ := http.NewRequest("POST", config.BaseURL+"/endpoint", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating endpoint:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to create endpoint. Status:", resp.Status)
		return
	}

	var createResponse struct {
		MockURL string `json:"mockUrl"`
	}
	json.NewDecoder(resp.Body).Decode(&createResponse)

	fmt.Println("Endpoint created successfully!")
	fmt.Println("Mock URL:", createResponse.MockURL)
}
