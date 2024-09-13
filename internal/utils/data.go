package utils

import (
	"encoding/json"
	"errors"
	"fmt"

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
	var result map[interface{}]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return convertMap(result), nil
}

func ParseYAML(data string) (map[string]interface{}, error) {
	if !IsYAML(data) {
		return nil, errors.New("data is not a valid YAML")
	}
	var result map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return convertMap(result), nil
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

func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convertMap(v)
		case []interface{}:
			res[fmt.Sprint(k)] = convertSlice(v)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}

func convertSlice(s []interface{}) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			res[i] = convertMap(v)
		case []interface{}:
			res[i] = convertSlice(v)
		default:
			res[i] = v
		}
	}
	return res
}
