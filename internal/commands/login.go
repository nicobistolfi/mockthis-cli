package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	fmt.Println("Login hash:", loginResponse.LoginHash)

	// In a real implementation, you'd wait for the user to click the magic link
	// and then exchange the login hash for a JWT. For this example, we'll just
	// save the email and a dummy token.
	config.SaveConfig(config.EmailFile, email)
	config.SaveConfig(config.TokenFile, "dummy_jwt_token")
}

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
