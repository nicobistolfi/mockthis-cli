package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

// RegisterCmd is the command to register a new user
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run:   register,
}

func register(cmd *cobra.Command, args []string) {
	fullName := promptForInput("Enter your full name: ")
	email := promptForInput("Enter your email: ")
	githubHandle := promptForInput("Enter your GitHub handle: ")
	country := promptForInput("Enter your country: ")

	registerData := map[string]string{
		"fullName":     fullName,
		"email":        email,
		"githubHandle": githubHandle,
		"country":      country,
	}
	jsonData, _ := json.Marshal(registerData)

	resp, err := http.Post(config.BaseURL+"/register", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error sending registration request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Registration failed. Please try again.")
		return
	}

	var registerResponse struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println(registerResponse.Message)

	// Call login command after successful registration
	LoginCmd.Run(cmd, []string{email})
}
