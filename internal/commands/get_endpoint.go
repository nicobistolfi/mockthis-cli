package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

// GetEndpointCmd is the command to get details of an existing mock endpoint
var GetEndpointCmd = &cobra.Command{
	Use:   "get [id or mockIdentifier]",
	Short: "Get details of an endpoint",
	Args:  cobra.ExactArgs(1),
	Run:   getEndpointCmd,
}

var outputFormat string

func init() {
	GetEndpointCmd.Flags().StringVarP(&outputFormat, "output", "o", "list", "Output format: list, table, or json")
}

func getEndpointCmd(cmd *cobra.Command, args []string) {
	configData, err := config.LoadConfig(config.TokenFile)
	if err != nil {
		fmt.Println("You need to login first.")
		return
	}
	token := configData.Token

	req, _ := http.NewRequest("GET", config.BaseURL+"/endpoints", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching endpoints:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch endpoints. Status:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var endpoints []map[string]interface{}
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	idOrMockIdentifier := args[0]
	var targetEndpoint map[string]interface{}

	for _, endpoint := range endpoints {
		if endpoint["id"] == idOrMockIdentifier || endpoint["mockIdentifier"] == idOrMockIdentifier {
			targetEndpoint = endpoint
			break
		}
	}

	if targetEndpoint == nil {
		fmt.Println("Endpoint not found.")
		return
	}

	switch outputFormat {
	case "list":
		printEndpointDetails(targetEndpoint)
	case "table":
		printEndpointTable(targetEndpoint)
	case "json":
		printEndpointJSON(targetEndpoint)
	default:
		fmt.Println("Invalid output format. Using default list format.")
		printEndpointDetails(targetEndpoint)
	}
}

func printEndpointDetails(endpoint map[string]interface{}) {
	fmt.Println("Endpoint Details:")
	fmt.Printf("ID: %s\n", endpoint["id"])
	fmt.Printf("Mock Identifier: %s\n", endpoint["mockIdentifier"])
	fmt.Printf("HTTP Status: %d\n", int(endpoint["httpStatus"].(float64)))
	fmt.Printf("Created At: %s\n", endpoint["createdAt"])
	fmt.Printf("Endpoint URL: %s\n", endpoint["endpointUrl"])
	fmt.Printf("Response Content Type: %s\n", endpoint["responseContentType"])
	fmt.Printf("Charset: %s\n", endpoint["charset"])
	fmt.Printf("Response Body: %s\n", endpoint["responseBody"])

	if httpHeaders, ok := endpoint["httpHeaders"].(map[string]interface{}); ok {
		fmt.Println("HTTP Headers:")
		for key, value := range httpHeaders {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	if authCredentials, ok := endpoint["authCredentials"].(map[string]interface{}); ok {
		fmt.Println("Auth Credentials:")
		fmt.Printf("  Type: %s\n", authCredentials["type"])
		fmt.Printf("  Token: %s\n", authCredentials["token"])
	}

	fmt.Printf("CURL: %s\n", endpoint["curl"])
}

func printEndpointTable(endpoint map[string]interface{}) {
	fmt.Println("| Key | Value |")
	fmt.Println("|-----|-------|")
	fmt.Printf("| ID | %s |\n", endpoint["id"])
	fmt.Printf("| Mock Identifier | %s |\n", endpoint["mockIdentifier"])
	fmt.Printf("| HTTP Status | %d |\n", int(endpoint["httpStatus"].(float64)))
	fmt.Printf("| Created At | %s |\n", endpoint["createdAt"])
	fmt.Printf("| Endpoint URL | %s |\n", endpoint["endpointUrl"])
	fmt.Printf("| Response Content Type | %s |\n", endpoint["responseContentType"])
	fmt.Printf("| Charset | %s |\n", endpoint["charset"])
	fmt.Printf("| Response Body | %s |\n", endpoint["responseBody"])
	// ... Add other fields as needed ...
}

func printEndpointJSON(endpoint map[string]interface{}) {
	jsonData, err := json.MarshalIndent(endpoint, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
