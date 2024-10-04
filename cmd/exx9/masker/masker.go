package masker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// maskingStrategy defines a function type for masking
type maskingStrategy func(interface{}) interface{}

// defaultMaskingStrategy uses a simple replacement with "****"
func defaultMaskingStrategy(_ interface{}) interface{} {
	return "****"
}

// maskConfig holds the masking configuration
type maskConfig struct {
	Keys         []string        // Keys to be masked
	MaskStrategy maskingStrategy // Masking strategy to apply
}

// selectMaskingStrategy returns the user-defined or default masking strategy
func selectMaskingStrategy(strategy []maskingStrategy) maskingStrategy {
	if len(strategy) > 0 {
		return strategy[0]
	}
	return defaultMaskingStrategy
}

// MaskByKeys masks values for specified keys
// Params:
// - jsonData: JSON string
// - keys: list of keys whose values should be masked
// - strategy: custom masking strategy (optional)
// Returns the masked JSON and an error (if any)
func MaskByKeys(jsonData string, keys []string, strategy ...maskingStrategy) (string, error) {
	maskStrategy := selectMaskingStrategy(strategy)
	config := maskConfig{
		Keys:         keys,
		MaskStrategy: maskStrategy,
	}

	// Call the masking function
	maskedJSON, err := maskJSON([]byte(jsonData), config)
	if err != nil {
		return "", err
	}
	return string(maskedJSON), nil
}

// maskJSON masks the JSON data according to the provided configuration
func maskJSON(input []byte, config maskConfig) ([]byte, error) {
	var data map[string]interface{}

	// Decode JSON into a map
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Mask the data
	maskedData, err := processMap(data, config, false, "")
	if err != nil {
		return nil, fmt.Errorf("error processing map: %w", err)
	}

	// Encode back into JSON
	maskedJSON, err := json.Marshal(maskedData)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	return maskedJSON, nil
}

// processMap processes the JSON object and masks values where necessary
func processMap(data map[string]interface{}, config maskConfig, maskAll bool, path string) (interface{}, error) {
	maskedData := make(map[string]interface{})
	var err error
	for key, value := range data {
		fullPath := appendPath(path, key)
		if contains(config.Keys, key, fullPath) || maskAll {
			// Mask the value
			maskedData[key], err = processValue(value, config, true, fullPath)
		} else {
			// Recursively process nested structures
			maskedData[key], err = processValue(value, config, false, fullPath)
		}
		if err != nil {
			return nil, err
		}
	}
	return maskedData, nil
}

// processSlice processes an array of JSON data
func processSlice(data []interface{}, config maskConfig, maskAll bool, indexStart int, indexEnd int, path string) ([]interface{}, error) {
	if indexStart < 0 || indexEnd > len(data) || indexStart > indexEnd {
		return nil, fmt.Errorf("invalid range")
	}
	maskedSlice := make([]interface{}, len(data))
	var err error
	for i, value := range data {
		if maskAll || indexStart <= i && i < indexEnd {
			maskedSlice[i], err = processValue(value, config, true, path)

		} else {
			maskedSlice[i], err = processValue(value, config, false, path)
		}
		if err != nil {
			return nil, err
		}
	}
	return maskedSlice, nil
}

// processValue processes the value based on its type
func processValue(value interface{}, config maskConfig, toHash bool, path string) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		maskedValue, err := processMap(v, config, toHash, path)
		if err != nil {
			return nil, err
		}
		return maskedValue, nil
	case []interface{}:
		name, err := isSliceTag(config.Keys, path)
		if err != nil {
			return nil, err
		}

		var indexStart, indexEnd int
		if name != "" {
			indexStart, indexEnd, err = parseRange(name, len(v))
			if err != nil {
				return nil, err
			}
		}

		maskedValue, err := processSlice(v, config, toHash, indexStart, indexEnd, path)
		if err != nil {
			return nil, err
		}

		return maskedValue, nil
	case float64, int, bool, nil:
		if toHash {
			return config.MaskStrategy(v), nil
		}
		return v, nil
	case string:
		if toHash {
			return config.MaskStrategy(v), nil
		}
		return v, nil
	default:
		return "nil", nil
	}
}

// appendPath appends a key to the current path, forming the full path
func appendPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "/" + key
}

// contains checks if a key or path exists in the list of keys
func contains(slice []string, key string, path string) bool {
	for _, v := range slice {
		if v == key || v == path {
			return true
		}
	}
	return false
}

// isSliceTag checks if the path corresponds to a slice tag (e.g., "friends[0:5]")
func isSliceTag(slice []string, path string) (string, error) {
	for _, key := range slice {
		openBracket := strings.Index(key, "[")
		closeBracket := strings.Index(key, "]")
		if openBracket == -1 && closeBracket == -1 {
			continue
		}
		name := key[:openBracket]
		if openBracket == -1 || closeBracket == -1 || openBracket > closeBracket {
			return "", fmt.Errorf("invalid slice format in key: %s", key)
		}

		if strings.Contains(path, name) {
			return key, nil
		}
	}
	return "", nil
}

// parseRange parses a slice range (e.g., "friends[0:5]" or "[2:]")
func parseRange(s string, length int) (int, int, error) {
	var startIndex, endIndex int
	var err error
	openBracket := strings.Index(s, "[")
	closeBracket := strings.Index(s, "]")

	rangePart := s[openBracket+1 : closeBracket]
	colonIndex := strings.Index(rangePart, ":")
	if colonIndex == -1 {
		startIndex, err = strconv.Atoi(rangePart)
		if err != nil {
			return 0, 0, err
		}
		return startIndex, startIndex + 1, nil
	}

	if startStr := rangePart[:colonIndex]; startStr == "" {
		startIndex = 0
	} else {
		startIndex, err = strconv.Atoi(startStr)
		if err != nil {
			return 0, 0, err
		}
	}

	if endStr := rangePart[colonIndex+1:]; endStr == "" {
		endIndex = length
	} else {
		endIndex, err = strconv.Atoi(endStr)
		if err != nil {
			return 0, 0, err
		}
	}

	return startIndex, endIndex, nil
}
