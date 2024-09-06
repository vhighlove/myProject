package masker

import (
	"encoding/json"
	"errors"
	"fmt"
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

// extractKeysFromJSONString extracts all keys from JSON in proper order
func extractKeysFromJSONString(jsonStr string) ([]string, error) {
	// Remove spaces for easier processing
	jsonStr = strings.ReplaceAll(jsonStr, " ", "")

	var keys []string
	length := len(jsonStr)

	for i := 0; i < length; i++ {
		// Look for key between " " and :
		if jsonStr[i] == '"' {
			start := i + 1
			end := strings.Index(jsonStr[start:], "\"")
			if end == -1 {
				return nil, errors.New("error: closing quote not found")
			}
			end += start

			// Check if there's a colon after the key
			if end+1 < length && jsonStr[end+1] == ':' {
				key := jsonStr[start:end]
				keys = append(keys, key)
				i = end + 1 // Skip the key and colon
			}
		}
	}
	return keys, nil
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
		return "", fmt.Errorf("error during masking by keys: %w", err)
	}
	return string(maskedJSON), nil
}

// MaskByIndexes masks values by the indexes of the keys
// Params:
// - jsonData: JSON string
// - indexes: indexes of the keys whose values should be masked
// - strategy: custom masking strategy (optional)
// Returns the masked JSON and an error (if any)
func MaskByIndexes(jsonData string, indexes []int, strategy ...maskingStrategy) (string, error) {
	maskStrategy := selectMaskingStrategy(strategy)
	keys, err := extractKeysFromJSONString(jsonData)
	if err != nil {
		return "", fmt.Errorf("error extracting keys: %w", err)
	}

	if len(keys) == 0 {
		return "", errors.New("no keys found")
	}

	var newKeys []string
	for _, index := range indexes {
		if index < 0 || index >= len(keys) {
			return "", fmt.Errorf("index %d is out of bounds for keys", index)
		}
		newKeys = append(newKeys, keys[index])
	}

	config := maskConfig{
		Keys:         newKeys,
		MaskStrategy: maskStrategy,
	}

	maskedJSON, err := maskJSON([]byte(jsonData), config)
	if err != nil {
		return "", fmt.Errorf("error during masking by indexes: %w", err)
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
	maskedData, err := processMap(data, config, false)
	if err != nil {
		return nil, err
	}

	// Encode back into JSON
	maskedJSON, err := json.Marshal(maskedData)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	return maskedJSON, nil
}

// processMap processes the JSON object and masks values where necessary
func processMap(data map[string]interface{}, config maskConfig, maskAll bool) (interface{}, error) {
	maskedData := make(map[string]interface{})
	for key, value := range data {
		if contains(config.Keys, key) || maskAll {
			// Mask the value
			maskedData[key] = processValueForMasking(value, config)
		} else {
			// Recursively process nested structures
			maskedData[key] = processValue(value, config)
		}
	}
	return maskedData, nil
}

// processSlice processes an array of JSON data
func processSlice(data []interface{}, config maskConfig, maskAll bool) ([]interface{}, error) {
	maskedSlice := make([]interface{}, len(data))
	for i, value := range data {
		if maskAll {
			maskedSlice[i] = processValueForMasking(value, config)
			continue
		}
		maskedSlice[i] = processValue(value, config)
	}
	return maskedSlice, nil
}

// processValue processes the value based on its type
func processValue(value interface{}, config maskConfig) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		maskedValue, _ := processMap(v, config, false)
		return maskedValue
	case []interface{}:
		maskedValue, _ := processSlice(v, config, false)
		return maskedValue
	case float64:
		return v
	case int:
		return v
	case string:
		return v
	case bool:
		return v
	case nil:
		return v
	default:
		return "nil"
	}
}

// processValueForMasking processes the value and applies the masking strategy
func processValueForMasking(value interface{}, config maskConfig) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		maskedValue, _ := processMap(v, config, true)
		return maskedValue
	case []interface{}:
		maskedValue, _ := processSlice(v, config, true)
		return maskedValue
	case float64:
		return config.MaskStrategy(v)
	case int:
		return config.MaskStrategy(v)
	case string:
		return config.MaskStrategy(v)
	case bool:
		return config.MaskStrategy(v)
	case nil:
		return config.MaskStrategy(v)
	default:
		return "nil"
	}
}

// contains checks if an item exists in a list
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

//package masker
//
//import (
//	"encoding/json"
//	"strings"
//)
//
//// maskingStrategy визначає функцію маскування
//type maskingStrategy func(interface{}) interface{}
//
//// defaultMaskingStrategy використовує просте заміщення значенням "****"
//func defaultMaskingStrategy(_ interface{}) interface{} {
//	return "****"
//}
//
//// maskConfig зберігає конфігурації для маскування
//type maskConfig struct {
//	Keys         []string        // Ключі для маскування
//	MaskStrategy maskingStrategy // Стратегія маскування
//}
//
//// selectMaskingStrategy повертає стратегію маскування: користувацьку або дефолтну
//func selectMaskingStrategy(strategy []maskingStrategy) maskingStrategy {
//	if len(strategy) > 0 {
//		return strategy[0]
//	}
//	return defaultMaskingStrategy
//}
//
//func extractKeysFromJSONString(jsonStr string) []string {
//	// Удаляем все пробелы для упрощенной обработки
//	jsonStr = strings.ReplaceAll(jsonStr, " ", "")
//
//	// Список для хранения ключей
//	var keys []string
//
//	// Индекс для обхода строки
//	length := len(jsonStr)
//	for i := 0; i < length; i++ {
//		// Находим начало строки с ключом, которое начинается с кавычки
//		if jsonStr[i] == '"' {
//			start := i + 1
//			// Ищем конец строки (следующую кавычку)
//			end := strings.Index(jsonStr[start:], "\"")
//			if end == -1 {
//				break // если нет закрывающей кавычки, выходим из цикла
//			}
//			end += start
//
//			// Проверяем, стоит ли за строкой двоеточие (:)
//			if end+1 < length && jsonStr[end+1] == ':' {
//				// Извлекаем ключ и добавляем его в список
//				key := jsonStr[start:end]
//				keys = append(keys, key)
//				// Сдвигаем индекс на конец обработанного ключа
//				i = end + 1
//			}
//		}
//	}
//
//	return keys
//}
//
//// MaskByKeys маскує значення за ключами
//func MaskByKeys(jsonData string, keys []string, strategy ...maskingStrategy) (string, error) {
//	maskStrategy := selectMaskingStrategy(strategy)
//
//	config := maskConfig{
//		Keys:         keys,
//		MaskStrategy: maskStrategy,
//	}
//
//	maskedJSON, err := maskJSON([]byte(jsonData), config)
//	return string(maskedJSON), err
//}
//
//// MaskByIndexes маскує значення за шляхами
//func MaskByIndexes(jsonData string, indexes []int, strategy ...maskingStrategy) (string, error) {
//	maskStrategy := selectMaskingStrategy(strategy)
//	keys := extractKeysFromJSONString(jsonData)
//	newKeys := []string{}
//	for _, index := range indexes {
//		newKeys = append(newKeys, keys[index])
//	}
//	config := maskConfig{
//		Keys:         newKeys,
//		MaskStrategy: maskStrategy,
//	}
//
//	maskedJSON, err := maskJSON([]byte(jsonData), config)
//	return string(maskedJSON), err
//}
//
//// maskJSON - загальна функція для обробки маскування
//func maskJSON(input []byte, config maskConfig) ([]byte, error) {
//	var data map[string]interface{}
//	if err := json.Unmarshal(input, &data); err != nil {
//		return nil, err
//	}
//
//	maskedData, err := processMap(data, config, false)
//	if err != nil {
//		return nil, err
//	}
//
//	return json.Marshal(maskedData)
//}
//
//// processValue обробляє значення залежно від типу
//func processValue(value interface{}, config maskConfig) interface{} {
//	switch v := value.(type) {
//	case map[string]interface{}:
//		maskedValue, _ := processMap(v, config, false)
//		return maskedValue
//	case []interface{}:
//		maskedValue, _ := processSlice(v, config, false)
//		return maskedValue
//	case float64:
//		return v
//	case int:
//		return v
//	case string:
//		return v
//	case bool:
//		return v
//	case nil:
//		return v
//	default:
//		return "nil"
//	}
//}
//
//// processMap обробляє JSON об'єкт (map), зберігаючи порядок ключів
//func processMap(data map[string]interface{}, config maskConfig, maskAll bool) (interface{}, error) {
//	maskedData := make(map[string]interface{})
//	for key, value := range data {
//		// Перевіряємо, чи ключ є у списку ключів для маскування
//		if contains(config.Keys, key) || maskAll {
//			// Якщо значення - простий тип, маскуємо його
//			maskedData[key] = processValueForMasking(value, config)
//		} else {
//			// Інакше, обробляємо значення в залежності від його типу
//			maskedData[key] = processValue(value, config)
//		}
//	}
//	return maskedData, nil
//}
//
//// processSlice обробляє JSON масив
//func processSlice(data []interface{}, config maskConfig, maskAll bool) ([]interface{}, error) {
//	maskedSlice := make([]interface{}, len(data))
//	for i, value := range data {
//		if maskAll {
//			maskedSlice[i] = processValueForMasking(value, config)
//			continue
//		}
//		maskedSlice[i] = processValue(value, config)
//	}
//	return maskedSlice, nil
//}
//
//// processValueForMasking обробляє значення залежно від типу, з урахуванням маскування
//func processValueForMasking(value interface{}, config maskConfig) interface{} {
//	switch v := value.(type) {
//	case map[string]interface{}:
//		// Рекурсивно обробляємо вкладені об'єкти
//		maskedValue, _ := processMap(v, config, true)
//		return maskedValue
//	case []interface{}:
//		// Рекурсивно обробляємо вкладені масиви
//		maskedValue, _ := processSlice(v, config, true)
//		return maskedValue
//	case float64:
//		return config.MaskStrategy(v)
//	case int:
//		return config.MaskStrategy(v)
//	case string:
//		return config.MaskStrategy(v)
//	case bool:
//		return config.MaskStrategy(v)
//	case nil:
//		return config.MaskStrategy(v)
//	default:
//		return "nil"
//	}
//}
//
//func contains(slice []string, item string) bool {
//	for _, v := range slice {
//		if v == item {
//			return true
//		}
//	}
//	return false
//}
