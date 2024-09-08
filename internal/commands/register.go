package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

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
		UserID  string `json:"userId"`
	}
	json.NewDecoder(resp.Body).Decode(&registerResponse)

	fmt.Println(registerResponse.Message)
	fmt.Println("User ID:", registerResponse.UserID)
}
