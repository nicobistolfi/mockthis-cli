package utils

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateAgainstSchema validates the given data against a JSON schema file
func ValidateAgainstSchema(data interface{}, schema string) error {
	// Load the schema
	schemaLoader := gojsonschema.NewStringLoader(schema)

	// Load the data
	dataLoader := gojsonschema.NewGoLoader(data)

	// Perform the validation
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return fmt.Errorf("error during schema validation: %v", err)
	}

	// Check if the validation was successful
	if !result.Valid() {
		// Collect all validation errors
		var errorMessages string
		for _, desc := range result.Errors() {
			errorMessages += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("the data is not valid according to the schema:\n%s", errorMessages)
	}

	return nil
}
