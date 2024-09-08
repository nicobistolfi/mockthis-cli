package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var endpoints []struct {
		ID             string    `json:"id"`
		HttpStatus     int       `json:"httpStatus"`
		CreatedAt      time.Time `json:"createdAt"`
		MockIdentifier string    `json:"mockIdentifier"`
		EndpointURL    string    `json:"endpointUrl"`
	}

	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mock Identifier\tID\tMethod\tStatus\tCreated At\tEndpoint URL")
	fmt.Fprintln(w, "----------------\t--\t------\t------\t----------\t------------")

	for _, e := range endpoints {
		method := "GET" // Default method
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
			e.MockIdentifier,
			e.ID,
			method,
			e.HttpStatus,
			e.CreatedAt.Format("2006-01-02 15:04:05"),
			e.EndpointURL)
	}

	w.Flush()
}
