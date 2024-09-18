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
	"github.com/nicobistolfi/mockthis-cli/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CreateEndpointCmd is the command to create a new mock endpoint
var CreateEndpointCmd = &cobra.Command{
	Use:   "create [--file <path>] [--auth-type <type>] [--auth-properties <properties>] [--request-content-type <type>] [--request-schema <schema>] [--method <method>] [--status <status>] [--content-type <type>] [--charset <charset>] [--headers <headers>] [--schema <schema>] [--body <body>]",
	Short: "Create a new mock endpoint",
	Run:   createEndpoint,
}

func init() {
	// File
	CreateEndpointCmd.Flags().StringP("file", "f", "", "Path to JSON or YAML file containing endpoint data")

	// Response
	CreateEndpointCmd.Flags().StringP("method", "m", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	CreateEndpointCmd.Flags().StringP("status", "s", "200", "HTTP status code")
	CreateEndpointCmd.Flags().StringP("content-type", "c", "application/json", "Response Content-Type")
	CreateEndpointCmd.Flags().String("charset", "", "Charset")
	CreateEndpointCmd.Flags().StringP("headers", "H", "", "Response headers (comma-separated key=value pairs)")
	CreateEndpointCmd.Flags().String("schema", "", "JSON Schema to validate the response body")
	CreateEndpointCmd.Flags().StringP("body", "b", "Hello, World! 🌎", "Response body")

	// Authentication
	CreateEndpointCmd.Flags().String("auth-type", "", "Authentication type (basic, apiKey, bearer, oauth2, jwt)")
	CreateEndpointCmd.Flags().String("auth-properties", "", "Authentication properties (comma-separated key=value pairs)")

	// Request
	CreateEndpointCmd.Flags().String("request-content-type", "application/json", "Request Content-Type")
	CreateEndpointCmd.Flags().String("request-schema", "", "JSON Schema to validate the request body")
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
		err := loadFromFile(filePath, cmd)
		if err != nil {
			return nil, err
		}
	}
	endpointData = loadFromFlags(cmd)

	// Convert status to int or set to 200 if not present or conversion fails
	if status, ok := endpointData["status"].(string); ok {
		if convertedStatus, err := strconv.Atoi(status); err == nil {
			endpointData["status"] = convertedStatus
		} else {
			endpointData["status"] = 200
		}
	} else {
		endpointData["status"] = 200
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
		endpointData["responseBody"] = nil
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

	req, _ := http.NewRequest("POST", config.BaseURL+"/endpoints", bytes.NewBuffer(jsonData))
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
		MockURL  string                 `json:"mockUrl"`
		ID       string                 `json:"id"`
		Endpoint map[string]interface{} `json:"endpoint"`
	}
	err := json.NewDecoder(resp.Body).Decode(&createResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding API response: %v", err)
	}

	table := buildTableFromMap(createResponse.Endpoint)

	return fmt.Sprintf("Endpoint created successfully!\nMock URL: %s\n\n%s", createResponse.MockURL, table), nil
}

// Print a table from a map and keep it aligned to the left
func buildTableFromMap(data map[string]interface{}) string {
	// Define the order of keys
	order := []string{"ID", "Method", "Status", "ResponseContentType", "ResponseBody"}

	// Find the maximum width for each column
	maxKeyWidth := 5   // minimum width for "Field"
	maxValueWidth := 5 // minimum width for "Value"
	for key, value := range data {
		keyWidth := len(key)
		valueWidth := len(fmt.Sprintf("%v", value))
		if keyWidth > maxKeyWidth {
			maxKeyWidth = keyWidth
		}
		if valueWidth > maxValueWidth {
			maxValueWidth = valueWidth
		}
	}

	// Create the format string for each row
	rowFormat := fmt.Sprintf("| %%-%ds | %%-%ds |\n", maxKeyWidth, maxValueWidth)

	// Build the table
	var tableBuilder strings.Builder
	tableBuilder.WriteString(fmt.Sprintf(rowFormat, "Field", "Value"))
	tableBuilder.WriteString(fmt.Sprintf("+%s+%s+\n", strings.Repeat("-", maxKeyWidth+2), strings.Repeat("-", maxValueWidth+2)))

	// First, add rows for ordered keys
	for _, key := range order {
		if value, exists := data[key]; exists && value != nil && value != "" {
			fmt.Fprintf(&tableBuilder, rowFormat, key, formatValue(value))
			delete(data, key)
		}
	}

	// Then, add rows for remaining keys
	for key, value := range data {
		if value != nil && value != "" {
			fmt.Fprintf(&tableBuilder, rowFormat, key, formatValue(value))
		}
	}

	return tableBuilder.String()
}

func formatValue(value interface{}) string {
	if m, ok := value.(map[string]interface{}); ok {
		pairs := make([]string, 0, len(m))
		for k, v := range m {
			if v != nil && v != "" {
				pairs = append(pairs, fmt.Sprintf("%s: %v", k, v))
			}
		}
		if len(pairs) > 0 {
			return "[ " + strings.Join(pairs, ", ") + " ]"
		}
		return "[]"
	}
	return fmt.Sprintf("%v", value)
}

// ProcessAuthCredentials processes the authentication credentials from a string that can be either a JSON or a comma-separated list of key=value pairs
func processAuthCredentials(authType, authProperties string) map[string]interface{} {
	authCredentials := map[string]interface{}{
		"type": authType,
	}

	authPropertiesMap := make(map[string]interface{})
	if utils.IsJSON(authProperties) {
		err := json.Unmarshal([]byte(authProperties), &authPropertiesMap)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			os.Exit(1)
		}
	} else {
		propertyPairs := strings.Split(authProperties, ",")

		for _, pair := range propertyPairs {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				authPropertiesMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	switch authType {
	case "basic":
		authCredentials["username"] = authPropertiesMap["username"]
		authCredentials["password"] = authPropertiesMap["password"]
	case "apiKey":
		authCredentials["name"] = authPropertiesMap["name"]
		authCredentials["value"] = authPropertiesMap["value"]
		authCredentials["in"] = authPropertiesMap["in"]
	case "bearer":
		authCredentials["token"] = authPropertiesMap["token"]
	case "oauth2":
		authCredentials["accessToken"] = authPropertiesMap["accessToken"]
		authCredentials["tokenType"] = authPropertiesMap["tokenType"]
		authCredentials["expiresIn"] = authPropertiesMap["expiresIn"]
		authCredentials["refreshToken"] = authPropertiesMap["refreshToken"]
	case "jwt":
		authCredentials["token"] = authPropertiesMap["token"]
	default:
		fmt.Println("Invalid authentication type. Supported types: basic, apikey, bearer, oauth2, jwt. Got:", authType)
		os.Exit(1)
	}

	return authCredentials
}

func loadFromFile(filePath string, cmd *cobra.Command) error {
	data, err := utils.LoadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}

	var endpointData map[string]interface{}
	ext := filepath.Ext(filePath)

	switch {
	case utils.IsJSON(data):
		endpointData, err = utils.ParseJSON(data)
	case utils.IsYAML(data):
		endpointData, err = utils.ParseYAML(data)
	default:
		fmt.Printf("Unsupported file format: %s. Use JSON or YAML.\n", ext)
		err = fmt.Errorf("unsupported file format: %s", ext)
		return err
	}

	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return err
	}

	// Check if "endpoint" key exists and is a map
	endpoint, ok := endpointData["endpoint"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("endpoint data is not a map or the endpoint key is not found")
	}

	utils.MapToFlags(endpoint, cmd)
	return nil
}

func loadFromFlags(cmd *cobra.Command) map[string]interface{} {
	endpointData := make(map[string]interface{})

	// Custom flag mappings
	flagMappings := map[string]string{
		"body":                 "responseBody",
		"content-type":         "responseContentType",
		"request-content-type": "requestContentType",
		"schema":               "responseBodySchema",
		"request-schema":       "requestBodySchema",
		"headers":              "httpHeaders",
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
