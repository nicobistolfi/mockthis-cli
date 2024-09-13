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
    http-status: 200
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
	createFlags(cmd, data["endpoint"].(map[string]interface{}))

	MapToFlags(data["endpoint"].(map[string]interface{}), cmd)

	expectedFlags := []struct {
		name  string
		value string
	}{
		{"auth-type", "basic"},
		{"auth-properties-username", "admin"},
		{"auth-properties-password", "admin"},
		{"response-method", "GET"},
		{"response-http-status", "200"},
		{"response-content-type", "application/json"},
		{"response-charset", "UTF-8"},
		{"response-headers-Content-Type", "application/json"},
		{"response-headers-X-Random-Header", "MockThis Random Header"},
		{"response-schema-type", "string"},
		{"response-body", "Hello, World! 🌎"},
		{"request-content-type", "application/json"},
		{"request-schema-type", "string"},
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

	// Create flags before calling MapToFlags
	createFlags(cmd, data["endpoint"].(map[string]interface{}))

	MapToFlags(data["endpoint"].(map[string]interface{}), cmd)

	expectedFlags := []struct {
		name  string
		value string
	}{
		{"response-body", "{\n  \"hello\": \"world\",\n  \"foo\": \"bar\"\n}\n"},
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
