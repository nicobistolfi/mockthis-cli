package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nicobistolfi/mockthis-cli/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var CreateEndpointCmd = &cobra.Command{
	Use:   "create [--file <path>] [--method <method>] [--http-status <code>] [--content-type <type>] [--charset <charset>] [--body <body>] [--auth-type <type>] [--auth-properties <properties>]",
	Short: "Create a new mock endpoint",
	Run:   createEndpoint,
}

func init() {
	CreateEndpointCmd.Flags().StringP("file", "f", "", "Path to JSON or YAML file containing endpoint data")
	CreateEndpointCmd.Flags().String("method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	CreateEndpointCmd.Flags().String("http-status", "200", "HTTP status code")
	CreateEndpointCmd.Flags().String("content-type", "application/json", "Response Content-Type")
	CreateEndpointCmd.Flags().String("request-content-type", "application/json", "Request Content-Type")
	CreateEndpointCmd.Flags().String("charset", "UTF-8", "Charset")
	CreateEndpointCmd.Flags().String("body", "Hello, World! 🌎", "Response body")
	CreateEndpointCmd.Flags().String("response-body-schema", "", "JSON Schema to validate the response body")
	CreateEndpointCmd.Flags().String("request-body-schema", "", "JSON Schema to validate the request body")
	CreateEndpointCmd.Flags().String("auth-type", "", "Authentication type (basic, api-key, bearer-token, oauth2, jwt)")
	CreateEndpointCmd.Flags().String("auth-properties", "", "Authentication properties (comma-separated key=value pairs)")
}

func createEndpoint(cmd *cobra.Command, args []string) {
	configData, err := config.LoadConfig(config.TokenFile)
	if err != nil {
		fmt.Println("You need to login first.")
		return
	}
	token := configData.Token

	var endpointData map[string]interface{}

	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		endpointData = loadFromFile(filePath)
	} else {
		endpointData = loadFromFlags(cmd)
	}

	// Convert httpStatus to int
	if httpStatus, ok := endpointData["httpStatus"].(string); ok {
		endpointData["httpStatus"], _ = strconv.Atoi(httpStatus)
	}

	// Add default values for httpHeaders and mockIdentifier
	if _, ok := endpointData["httpHeaders"]; !ok {
		endpointData["httpHeaders"] = map[string]string{}
	}
	if _, ok := endpointData["mockIdentifier"]; !ok {
		endpointData["mockIdentifier"] = "default-mock-endpoint"
	}

	authType, _ := cmd.Flags().GetString("auth-type")
	authProperties, _ := cmd.Flags().GetString("auth-properties")

	if authType != "" && authProperties != "" {
		endpointData["authCredentials"] = processAuthCredentials(authType, authProperties)
	} else if (authType != "" && authProperties == "") || (authType == "" && authProperties != "") {
		fmt.Println("auth-type and auth-properties are required when using authentication.")
		os.Exit(1)
	}

	endpointData["method"] = endpointData["method"].(string)
	endpointData["responseContentType"] = endpointData["responseContentType"].(string)
	endpointData["charset"] = endpointData["charset"].(string)
	endpointData["responseBody"] = endpointData["responseBody"].(string)

	if _, ok := endpointData["responseBodySchema"]; ok {
		endpointData["responseBodySchema"] = endpointData["responseBodySchema"].(string)
	}

	if _, ok := endpointData["requestContentType"]; ok {
		endpointData["requestContentType"] = endpointData["requestContentType"].(string)
	}

	if _, ok := endpointData["requestBodySchema"]; ok {
		endpointData["requestBodySchema"] = endpointData["requestBodySchema"].(string)
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

// New private function
func processAuthCredentials(authType, authProperties string) map[string]interface{} {
	authCredentials := map[string]interface{}{
		"type": authType,
	}

	authPropertiesMap := make(map[string]interface{})
	propertyPairs := strings.Split(authProperties, ",")

	for _, pair := range propertyPairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			authPropertiesMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	switch authType {
	case "basic":
		authCredentials["username"] = authPropertiesMap["username"]
		authCredentials["password"] = authPropertiesMap["password"]
	case "api-key":
		authCredentials["name"] = authPropertiesMap["name"]
		authCredentials["value"] = authPropertiesMap["value"]
		authCredentials["in"] = authPropertiesMap["in"]
	case "bearer-token":
		authCredentials["token"] = authPropertiesMap["token"]
	case "oauth2":
		authCredentials["accessToken"] = authPropertiesMap["accessToken"]
		authCredentials["tokenType"] = authPropertiesMap["tokenType"]
		authCredentials["expiresIn"] = authPropertiesMap["expiresIn"]
		authCredentials["refreshToken"] = authPropertiesMap["refreshToken"]
	case "jwt":
		authCredentials["token"] = authPropertiesMap["token"]
	default:
		fmt.Println("Invalid authentication type. Supported types: basic, api-key, bearer-token, oauth2, jwt")
		os.Exit(1)
	}

	return authCredentials
}

func loadFromFile(filePath string) map[string]interface{} {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	var endpointData map[string]interface{}
	if ext := filepath.Ext(filePath); ext == ".json" {
		err = json.Unmarshal(data, &endpointData)
	} else if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(data, &endpointData)
	} else {
		fmt.Println("Unsupported file format. Use JSON or YAML.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	return endpointData
}

func loadFromFlags(cmd *cobra.Command) map[string]interface{} {
	method, _ := cmd.Flags().GetString("method")
	httpStatus, _ := cmd.Flags().GetString("http-status")
	contentType, _ := cmd.Flags().GetString("content-type")
	charset, _ := cmd.Flags().GetString("charset")
	body, _ := cmd.Flags().GetString("body")

	if httpStatus == "" || contentType == "" || charset == "" || body == "" {
		fmt.Println("All flags (--method, --http-status, --content-type, --charset, --body) are required when not using a file.")
		os.Exit(1)
	}

	return map[string]interface{}{
		"method":              method,
		"httpStatus":          httpStatus,
		"responseContentType": contentType,
		"charset":             charset,
		"responseBody":        body,
	}
}
