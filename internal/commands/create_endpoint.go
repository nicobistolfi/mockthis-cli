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
	"github.com/spf13/pflag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// CreateEndpointCmd is the command to create a new mock endpoint
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
	endpointData, err := parseCommandArguments(cmd)
	if err != nil {
		fmt.Println("Error parsing command arguments:", err)
		os.Exit(1)
	}
	response := queryAPIEndpoint(endpointData)
	message, err := processAPIResponse(response)
	if err != nil {
		fmt.Println("Error processing API response:", err)
		os.Exit(1)
	}
	fmt.Println(message)
}

func parseCommandArguments(cmd *cobra.Command) (map[string]interface{}, error) {

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

	authType, _ := cmd.Flags().GetString("auth-type")
	authProperties, _ := cmd.Flags().GetString("auth-properties")

	if authType != "" && authProperties != "" {
		endpointData["authCredentials"] = processAuthCredentials(authType, authProperties)
	} else if (authType != "" && authProperties == "") || (authType == "" && authProperties != "") {
		fmt.Println("auth-type and auth-properties are required when using authentication.")
		os.Exit(1)
	}

	// Ensure correct types for specific fields
	if method, ok := endpointData["method"].(string); ok {
		endpointData["method"] = method
	} else {
		endpointData["method"] = "GET"
	}

	if responseContentType, ok := endpointData["responseContentType"].(string); ok {
		endpointData["responseContentType"] = responseContentType
	} else {
		endpointData["responseContentType"] = "application/json"
	}

	if charset, ok := endpointData["charset"].(string); ok {
		endpointData["charset"] = charset
	} else {
		endpointData["charset"] = "UTF-8"
	}

	if responseBody, ok := endpointData["responseBody"].(string); ok {
		endpointData["responseBody"] = responseBody
	} else {
		return nil, fmt.Errorf("responseBody is not a string or is missing")
	}

	for _, field := range []string{"responseBodySchema", "requestContentType", "requestBodySchema"} {
		if value, ok := endpointData[field].(string); ok {
			endpointData[field] = value
		} else {
			// If the field is missing or not a string, remove it from the map
			delete(endpointData, field)
		}
	}

	return endpointData, nil
}

func getConfig() *config.Data {
	configData, err := config.LoadConfig(config.TokenFile)
	if err != nil {
		fmt.Println("You need to login first.")
		os.Exit(1)
	}
	return configData
}

func queryAPIEndpoint(endpointData map[string]interface{}) *http.Response {
	configData := getConfig()

	jsonData, _ := json.Marshal(endpointData)

	req, _ := http.NewRequest("POST", config.BaseURL+"/endpoint", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+configData.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating endpoint:", err)
		os.Exit(1)
	}

	return resp
}

func processAPIResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create endpoint. Status: %s", resp.Status)
	}

	var createResponse struct {
		MockURL string `json:"mockUrl"`
	}
	err := json.NewDecoder(resp.Body).Decode(&createResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding API response: %v", err)
	}

	return fmt.Sprintf("Endpoint created successfully!\nMock URL: %s", createResponse.MockURL), nil
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
	endpointData := make(map[string]interface{})

	// Custom flag mappings
	flagMappings := map[string]string{
		"body":                 "responseBody",
		"content-type":         "responseContentType",
		"request-content-type": "requestContentType",
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Changed {
			key := flag.Name
			if mappedKey, exists := flagMappings[key]; exists {
				key = mappedKey
			} else {
				key = toCamelCase(key)
			}
			value := flag.Value.String()
			endpointData[key] = value
		}
	})

	return endpointData
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "-")
	caser := cases.Title(language.Und)
	for i := 1; i < len(parts); i++ {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, "")
}
