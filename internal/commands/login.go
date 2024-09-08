package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to MockThis",
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	email := promptForInput("Enter your email: ")

	loginData := map[string]string{"email": email}
	jsonData, _ := json.Marshal(loginData)

	resp, err := http.Post(config.BaseURL+"/login", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error sending login request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Login failed. Please try again.")
		return
	}

	var loginResponse struct {
		Message   string `json:"message"`
		LoginHash string `json:"loginHash"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResponse)

	fmt.Println(loginResponse.Message)
	fmt.Println("Please check your email and click the magic link.")

	// Poll for token
	token := pollForToken(email, loginResponse.LoginHash)
	if token == "" {
		fmt.Println("Login failed. Please try again.")
		return
	}

	// Save the token
	saveCredentials(email, token)
	fmt.Println("Login successful!")
}

func pollForToken(email, hash string) string {
	client := &http.Client{}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			url := fmt.Sprintf("%s/login/hash?email=%s&hash=%s", config.BaseURL, email, hash)
			resp, err := client.Get(url)
			if err != nil {
				fmt.Println("Error checking login status:", err)
				continue
			}
			defer resp.Body.Close()

			var response struct {
				LoginHashVerified bool   `json:"login_hash_verified"`
				Token             string `json:"token"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				fmt.Println("Error decoding response:", err)
				continue
			}

			if response.LoginHashVerified {
				return response.Token
			}
		case <-time.After(5 * time.Minute):
			fmt.Println("Login timed out. Please try again.")
			return ""
		}
	}
}

func saveCredentials(email, token string) {
	credentials := struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Email: email,
		Token: token,
	}

	jsonData, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		fmt.Println("Error creating credentials JSON:", err)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	credentialsPath := filepath.Join(homeDir, ".mockthis", "credentials")
	if err := os.MkdirAll(filepath.Dir(credentialsPath), 0700); err != nil {
		fmt.Println("Error creating .mockthis directory:", err)
		return
	}

	if err := os.WriteFile(credentialsPath, jsonData, 0600); err != nil {
		fmt.Println("Error saving credentials:", err)
		return
	}
}

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
