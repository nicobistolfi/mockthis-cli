package utils

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// MapToFlags sets values for existing flags in a cobra.Command
func MapToFlags(data map[string]interface{}, cmd *cobra.Command) {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			handleNestedMap(key, v, cmd)
		case string:
			cmd.Flags().Set(key, v)
		case int:
			cmd.Flags().Set(key, fmt.Sprintf("%d", v))
		case float64:
			cmd.Flags().Set(key, fmt.Sprintf("%f", v))
		case bool:
			cmd.Flags().Set(key, fmt.Sprintf("%t", v))
		}
	}

	// Print the flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("Flag: %s=%s\n", flag.Name, flag.Value)
	})
}

func handleNestedMap(prefix string, data map[string]interface{}, cmd *cobra.Command) {
	for key, value := range data {
		var fullKey string
		if prefix == "response" {
			fullKey = key
		} else {
			fullKey = prefix + "-" + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			// Stringify nested map instead of recursing
			jsonString, err := json.Marshal(v)
			if err == nil {
				cmd.Flags().Set(fullKey, string(jsonString))
			}
		case string:
			cmd.Flags().Set(fullKey, v)
		case int:
			cmd.Flags().Set(fullKey, fmt.Sprintf("%d", v))
		case float64:
			cmd.Flags().Set(fullKey, fmt.Sprintf("%f", v))
		case bool:
			cmd.Flags().Set(fullKey, fmt.Sprintf("%t", v))
		default:
			// Handle other types by converting to string
			cmd.Flags().Set(fullKey, fmt.Sprintf("%v", v))
		}
	}
}
