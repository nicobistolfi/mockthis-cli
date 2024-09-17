package commands

import (
	"bytes"
	"io"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			authType:       "apiKey",
			authProperties: "name=api_key,value=12345,in=header",
			expected: map[string]interface{}{
				"type":  "apiKey",
				"name":  "api_key",
				"value": "12345",
				"in":    "header",
			},
		},
		{
			name:           "Bearer Token",
			authType:       "bearer",
			authProperties: "token=abcdef123456",
			expected: map[string]interface{}{
				"type":  "bearer",
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

// nolint
func TestLoadFromFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("method", "", "")
	cmd.Flags().String("status", "", "")
	cmd.Flags().String("content-type", "", "")
	cmd.Flags().String("charset", "", "")
	cmd.Flags().String("body", "", "")

	cmd.Flags().Set("method", "GET")
	cmd.Flags().Set("status", "200")
	cmd.Flags().Set("content-type", "application/json")
	cmd.Flags().Set("charset", "UTF-8")
	cmd.Flags().Set("body", "Hello, World!")

	expected := map[string]interface{}{
		"method":              "GET",
		"status":              "200",
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
				"status":       "200",
				"content-type": "application/json",
				"charset":      "UTF-8",
				"body":         "Hello, World!",
			},
			expected: map[string]interface{}{
				"method":              "GET",
				"status":              200,
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
				//nolint
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
				Body:       io.NopCloser(bytes.NewBufferString(`{"MockURL": "https://dev.api.mockthis.io/api/v1/endpoints/1234567890"}`)),
			},
			expected: "Endpoint created successfully!\nMock URL: https://dev.api.mockthis.io/api/v1/endpoints/1234567890",
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
				if !strings.Contains(err.Error(), tt.expected) {
					t.Errorf("processAPIResponse() returned an error: %v, want it to contain %v", err, tt.expected)
				}
			} else {
				if !strings.Contains(result, tt.expected) {
					t.Errorf("processAPIResponse() = %v, want it to contain %v", result, tt.expected)
				}
			}

		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Get the directory of the current file
	_, filename, _, _ := runtime.Caller(0)
	currentDir := path.Dir(filename)
	testDataDir := path.Join(currentDir, "..", "..", "tests/data")

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
		expected map[string]string
	}{
		{
			name:     "Valid JSON file",
			filePath: path.Join(testDataDir, "valid.json"),
			wantErr:  false,
			expected: map[string]string{
				"method":       "POST",
				"status":       "201",
				"content-type": "application/json",
				"body":         `{"message":"Created"}`,
			},
		},
		{
			name:     "Valid YAML file",
			filePath: path.Join(testDataDir, "valid.yaml"),
			wantErr:  false,
			expected: map[string]string{
				"method":       "GET",
				"status":       "200",
				"content-type": "text/plain",
				"body":         "Hello, World!",
			},
		},
		{
			name:     "Invalid file format",
			filePath: path.Join(testDataDir, "invalid.txt"),
			wantErr:  true,
			expected: nil,
		},
		{
			name:     "Non-existent file",
			filePath: path.Join(testDataDir, "nonexistent.json"),
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			CreateEndpointCmd.Flags().VisitAll(func(flag *pflag.Flag) {
				cmd.Flags().AddFlag(flag)
			})

			err := loadFromFile(tt.filePath, cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for key, expectedValue := range tt.expected {
					flag := cmd.Flags().Lookup(key)
					if flag == nil {
						t.Errorf("Flag %s not found", key)
						continue
					}
					if flag.Value.String() != expectedValue {
						t.Errorf("Flag %s = %v, want %v", key, flag.Value.String(), expectedValue)
					}
				}
			}
		})
	}
}
