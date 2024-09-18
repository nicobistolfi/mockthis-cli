package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/truemail-rb/truemail-go"
)

// RegisterCmd is the command to register a new user
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run:   register,
}

func validateGithubHandle(githubHandle string) bool {
	url := fmt.Sprintf("https://api.github.com/users/%s", githubHandle)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error checking GitHub handle:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func register(cmd *cobra.Command, args []string) {
	fullName := promptForInput("Enter your full name: ")
	for len(fullName) <= 2 {
		fmt.Println("Full name must be at least 3 characters long. Please try again.")
		fullName = promptForInput("Enter your full name: ")
	}

	configuration, err := truemail.NewConfiguration(truemail.ConfigurationAttr{VerifierEmail: "nico@mockthis.io", ValidationTypeDefault: "regex"})
	if err != nil {
		fmt.Println("Error creating Truemail configuration:", err)
		return
	}
	var email string
	for {
		email = promptForInput("Enter your email: ")
		valid := truemail.IsValid(email, configuration)
		if !valid {
			fmt.Println("Invalid email format or unable to verify. Please try again.")
			continue
		}
		break
	}

	githubHandle := promptForInput("Enter your GitHub handle: ")
	for len(githubHandle) <= 2 || !validateGithubHandle(githubHandle) {
		if len(githubHandle) <= 2 {
			fmt.Println("GitHub handle must be at least 3 characters long. Please try again.")
		} else {
			fmt.Println("GitHub handle not found. Please enter a valid GitHub handle.")
		}
		githubHandle = promptForInput("Enter your GitHub handle: ")
	}

	country := promptForInput("Enter your country (2-letter code): ")
	country = strings.ToUpper(strings.TrimSpace(country))
	for len(country) != 2 {
		fmt.Println("Invalid country code. Please enter a 2-letter country code.")
		country = promptForInput("Enter your country (2-letter code): ")
		country = strings.ToUpper(strings.TrimSpace(country))
	}

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
