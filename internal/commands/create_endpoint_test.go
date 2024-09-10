package commands

import (
	"bytes"
	"io"
	"net/http"
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

func TestParseCommandArguments(t *testing.T) {
	tests := []struct {
		name     string
		flags    map[string]string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "Valid flags",
			flags: map[string]string{
				"method":       "GET",
				"http-status":  "200",
				"content-type": "application/json",
				"charset":      "UTF-8",
				"body":         "Hello, World!",
			},
			expected: map[string]interface{}{
				"method":              "GET",
				"httpStatus":          200,
				"responseContentType": "application/json",
				"charset":             "UTF-8",
				"responseBody":        "Hello, World!",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			for key, value := range tt.flags {
				cmd.Flags().String(key, "", "")
				cmd.Flags().Set(key, value)
			}
			result, err := parseCommandArguments(cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommandArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result == nil && tt.expected != nil {
				t.Errorf("parseCommandArguments() = nil, want %v", tt.expected)
				return
			}
			if result != nil {
				for key, expectedValue := range tt.expected {
					if resultValue, ok := result[key]; !ok {
						t.Errorf("parseCommandArguments() missing key %s", key)
					} else if resultValue != expectedValue {
						t.Errorf("parseCommandArguments() for key %s = %v, want %v", key, resultValue, expectedValue)
					}
				}
				for key := range result {
					if _, ok := tt.expected[key]; !ok {
						t.Errorf("parseCommandArguments() unexpected key %s", key)
					}
				}
			}
		})
	}
}

func TestProcessAPIResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *http.Response
		expected string
	}{
		{
			name: "Valid input",
			input: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"MockURL": "https://dev.api.mockthis.io/api/v1/endpoint/1234567890"}`)),
			},
			expected: "Endpoint created successfully!\nMock URL: https://dev.api.mockthis.io/api/v1/endpoint/1234567890",
		},
		{
			name: "Unauthorized",
			input: &http.Response{
				StatusCode: 401,
				Status:     "401 Unauthorized",
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Unauthorized"}`)),
			},
			expected: "failed to create endpoint. Status: 401 Unauthorized",
		},
		{
			name: "Invalid input",
			input: &http.Response{
				StatusCode: 500,
				Status:     "500 Internal Server Error",
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Internal Server Error"}`)),
			},
			expected: "failed to create endpoint. Status: 500 Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processAPIResponse(tt.input)
			if err != nil {
				// Handle the error
				if err.Error() != tt.expected {
					t.Errorf("processAPIResponse() returned an error: %v, want %v", err, tt.expected)
				}
			} else {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("processAPIResponse() = %v, want %v", result, tt.expected)
				}
			}

		})
	}
}
