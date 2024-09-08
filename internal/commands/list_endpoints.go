package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

var ListEndpointsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all created mock endpoints",
	Run:   listEndpoints,
}

func listEndpoints(cmd *cobra.Command, args []string) {
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
		fmt.Println("Error listing endpoints:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to list endpoints. Status:", resp.Status)
		return
	}

	var endpoints []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&endpoints)

	fmt.Println("Your mock endpoints:")
	for _, endpoint := range endpoints {
		fmt.Printf("ID: %s, URL: %s\n", endpoint["mockIdentifier"], endpoint["mockUrl"])
	}
}
