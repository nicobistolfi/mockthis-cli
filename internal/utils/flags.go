package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// MapToFlags sets values for existing flags in a cobra.Command
func MapToFlags(data map[string]interface{}, cmd *cobra.Command) {
	for key, value := range data {
		var err error
		switch v := value.(type) {
		case map[string]interface{}:
			handleNestedMap(key, v, cmd)
		case string:
			err = cmd.Flags().Set(key, v)
		case int:
			err = cmd.Flags().Set(key, fmt.Sprintf("%d", v))
		case float64:
			err = cmd.Flags().Set(key, fmt.Sprintf("%f", v))
		case bool:
			err = cmd.Flags().Set(key, fmt.Sprintf("%t", v))
		}
		if err != nil {
			fmt.Println("Error setting flag:", err)
			os.Exit(1)
		}
	}

	// Print the flags
	// cmd.Flags().VisitAll(func(flag *pflag.Flag) {
	// 	fmt.Printf("Flag: %s=%s\n", flag.Name, flag.Value)
	// })
}

func handleNestedMap(prefix string, data map[string]interface{}, cmd *cobra.Command) {
	for key, value := range data {
		var err error
		var fullKey string
		if prefix == "response" {
			fullKey = key
		} else {
			fullKey = prefix + "-" + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			// Stringify nested map instead of recursing
			jsonString, errJSON := json.Marshal(v)
			if errJSON != nil {
				fmt.Println("Error marshalling JSON:", errJSON)
				os.Exit(1)
			}
			err = cmd.Flags().Set(fullKey, string(jsonString))
		case string:
			err = cmd.Flags().Set(fullKey, v)
		case int:
			err = cmd.Flags().Set(fullKey, fmt.Sprintf("%d", v))
		case float64:
			err = cmd.Flags().Set(fullKey, fmt.Sprintf("%f", v))
		case bool:
			err = cmd.Flags().Set(fullKey, fmt.Sprintf("%t", v))
		default:
			// Handle other types by converting to string
			err = cmd.Flags().Set(fullKey, fmt.Sprintf("%v", v))
		}
		if err != nil {
			fmt.Println("Error setting flag:", err)
			os.Exit(1)
		}
	}
}
