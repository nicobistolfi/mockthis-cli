package utils

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"
)

func IsJSON(data string) bool {
	var result map[string]interface{}
	return json.Unmarshal([]byte(data), &result) == nil
}

func IsYAML(data string) bool {
	var result map[string]interface{}
	return yaml.Unmarshal([]byte(data), &result) == nil
}

func ParseJSON(data string) (map[string]interface{}, error) {
	if !IsJSON(data) {
		return nil, errors.New("data is not a valid JSON")
	}
	var result map[string]interface{}
	return result, json.Unmarshal([]byte(data), &result)
}

func ParseYAML(data string) (map[string]interface{}, error) {
	if !IsYAML(data) {
		return nil, errors.New("data is not a valid YAML")
	}
	var result map[string]interface{}
	return result, yaml.Unmarshal([]byte(data), &result)
}

func ToJSON(v interface{}) (string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func ToYAML(v interface{}) (string, error) {
	yamlBytes, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
