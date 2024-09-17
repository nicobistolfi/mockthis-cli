package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestMapToFlags(t *testing.T) {
	yamlInput := `
endpoint:
  auth:
    type: basic
    properties:
      username: admin
      password: admin
  response:
    method: GET
    status: 200
    content-type: application/json
    charset: UTF-8
    headers:
      Content-Type: application/json
      X-Random-Header: MockThis Random Header
    schema:
      type: string
    body: Hello, World! 🌎
  request:
    content-type: application/json
    schema:
      type: string
`

	data, err := ParseYAML(yamlInput)
	assert.NoError(t, err)

	cmd := &cobra.Command{}

	// Create flags before calling MapToFlags
	// File
	cmd.Flags().StringP("file", "f", "", "Path to JSON or YAML file containing endpoint data")

	// Response
	cmd.Flags().String("method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	cmd.Flags().String("status", "200", "HTTP status code")
	cmd.Flags().String("content-type", "application/json", "Response Content-Type")
	cmd.Flags().String("charset", "UTF-8", "Charset")
	cmd.Flags().String("headers", "", "Response headers (comma-separated key=value pairs)")
	cmd.Flags().String("schema", "", "JSON Schema to validate the response body")
	cmd.Flags().String("body", "Hello, World! 🌎", "Response body")

	// Authentication
	cmd.Flags().String("auth-type", "", "Authentication type (basic, apiKey, bearer, oauth2, jwt)")
	cmd.Flags().String("auth-properties", "", "Authentication properties (comma-separated key=value pairs)")

	// Request
	cmd.Flags().String("request-content-type", "application/json", "Request Content-Type")
	cmd.Flags().String("request-schema", "", "JSON Schema to validate the request body")

	MapToFlags(data["endpoint"].(map[string]interface{}), cmd)

	expectedFlags := []struct {
		name  string
		value string
	}{
		{"auth-type", "basic"},
		{"auth-properties", "{\"password\":\"admin\",\"username\":\"admin\"}"},
		{"method", "GET"},
		{"status", "200"},
		{"content-type", "application/json"},
		{"charset", "UTF-8"},
		{"headers", "{\"Content-Type\":\"application/json\",\"X-Random-Header\":\"MockThis Random Header\"}"},
		{"schema", "{\"type\":\"string\"}"},
		{"body", "Hello, World! 🌎"},
		{"request-content-type", "application/json"},
	}

	for _, ef := range expectedFlags {
		flag := cmd.Flags().Lookup(ef.name)
		assert.NotNil(t, flag, "Flag %s not found", ef.name)
		assert.Equal(t, ef.value, flag.Value.String(), "Unexpected value for flag %s", ef.name)
	}
}

func TestMapToFlagsJSONBody(t *testing.T) {
	yamlInput := `
endpoint:
  response:
    body:
      |
      {
        "hello": "world",
        "foo": "bar"
      }
`

	data, err := ParseYAML(yamlInput)
	assert.NoError(t, err)

	cmd := &cobra.Command{}

	cmd.Flags().String("body", "", "Response body")

	// Create flags before calling MapToFlags
	createFlags(cmd, data["endpoint"].(map[string]interface{}))

	MapToFlags(data["endpoint"].(map[string]interface{}), cmd)

	expectedFlags := []struct {
		name  string
		value string
	}{
		{"body", "{\n  \"hello\": \"world\",\n  \"foo\": \"bar\"\n}\n"},
	}

	for _, ef := range expectedFlags {
		flag := cmd.Flags().Lookup(ef.name)
		assert.NotNil(t, flag, "Flag %s not found", ef.name)
		assert.Equal(t, ef.value, flag.Value.String(), "Unexpected value for flag %s", ef.name)
	}
}

// Helper function to create flags recursively
func createFlags(cmd *cobra.Command, data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			createNestedFlags(cmd, key, v)
		default:
			cmd.Flags().String(key, "", "")
		}
	}
}

func createNestedFlags(cmd *cobra.Command, prefix string, data map[string]interface{}) {
	for key, value := range data {
		fullKey := prefix + "-" + key
		switch v := value.(type) {
		case map[string]interface{}:
			createNestedFlags(cmd, fullKey, v)
		default:
			cmd.Flags().String(fullKey, "", "")
		}
	}
}
