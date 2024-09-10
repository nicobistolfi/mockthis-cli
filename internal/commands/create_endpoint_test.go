package commands

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestProcessAuthCredentials(t *testing.T) {
	tests := []struct {
		name           string
		authType       string
		authProperties string
		expected       map[string]interface{}
	}{
		{
			name:           "Basic Auth",
			authType:       "basic",
			authProperties: "username=user,password=pass",
			expected: map[string]interface{}{
				"type":     "basic",
				"username": "user",
				"password": "pass",
			},
		},
		{
			name:           "API Key",
			authType:       "api-key",
			authProperties: "name=api_key,value=12345,in=header",
			expected: map[string]interface{}{
				"type":  "api-key",
				"name":  "api_key",
				"value": "12345",
				"in":    "header",
			},
		},
		{
			name:           "Bearer Token",
			authType:       "bearer-token",
			authProperties: "token=abcdef123456",
			expected: map[string]interface{}{
				"type":  "bearer-token",
				"token": "abcdef123456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processAuthCredentials(tt.authType, tt.authProperties)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("processAuthCredentials() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary JSON file
	jsonContent := `{"method": "POST", "httpStatus": 201, "responseContentType": "application/json"}`
	jsonFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(jsonFile.Name())
	jsonFile.Write([]byte(jsonContent))
	jsonFile.Close()

	// Create a temporary YAML file
	yamlContent := "method: GET\nhttpStatus: 200\nresponseContentType: text/plain"
	yamlFile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(yamlFile.Name())
	yamlFile.Write([]byte(yamlContent))
	yamlFile.Close()

	tests := []struct {
		name     string
		filePath string
		expected map[string]interface{}
	}{
		{
			name:     "JSON File",
			filePath: jsonFile.Name(),
			expected: map[string]interface{}{
				"method":              "POST",
				"httpStatus":          float64(201),
				"responseContentType": "application/json",
			},
		},
		{
			name:     "YAML File",
			filePath: yamlFile.Name(),
			expected: map[string]interface{}{
				"method":              "GET",
				"httpStatus":          200,
				"responseContentType": "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loadFromFile(tt.filePath)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("loadFromFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadFromFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("method", "", "")
	cmd.Flags().String("http-status", "", "")
	cmd.Flags().String("content-type", "", "")
	cmd.Flags().String("charset", "", "")
	cmd.Flags().String("body", "", "")

	cmd.Flags().Set("method", "GET")
	cmd.Flags().Set("http-status", "200")
	cmd.Flags().Set("content-type", "application/json")
	cmd.Flags().Set("charset", "UTF-8")
	cmd.Flags().Set("body", "Hello, World!")

	expected := map[string]interface{}{
		"method":              "GET",
		"httpStatus":          "200",
		"responseContentType": "application/json",
		"charset":             "UTF-8",
		"responseBody":        "Hello, World!",
	}

	result := loadFromFlags(cmd)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("loadFromFlags() = %v, want %v", result, expected)
	}
}

func TestCreateEndpoint(t *testing.T) {
	// This test is more complex as it involves mocking HTTP requests and responses.
	// For simplicity, we'll just test the JSON marshaling of the endpoint data.

	cmd := &cobra.Command{}
	cmd.Flags().String("method", "POST", "")
	cmd.Flags().String("http-status", "201", "")
	cmd.Flags().String("content-type", "application/json", "")
	cmd.Flags().String("charset", "UTF-8", "")
	cmd.Flags().String("body", `{"message": "Created"}`, "")

	endpointData := loadFromFlags(cmd)
	endpointData["httpStatus"], _ = json.Number("201").Int64()
	endpointData["httpHeaders"] = map[string]string{}
	endpointData["mockIdentifier"] = "default-mock-endpoint"
	endpointData["body"] = endpointData["responseBody"] // Add this line

	jsonData, err := json.Marshal(endpointData)
	if err != nil {
		t.Fatalf("Failed to marshal endpoint data: %v", err)
	}

	expected := `{"body":"{\"message\": \"Created\"}","charset":"UTF-8","httpHeaders":{},"httpStatus":201,"method":"POST","mockIdentifier":"default-mock-endpoint","responseBody":"{\"message\": \"Created\"}","responseContentType":"application/json"}`

	if string(jsonData) != expected {
		t.Errorf("Marshaled JSON does not match expected.\nGot:  %s\nWant: %s", string(jsonData), expected)
	}
}
