package mapping

import (
	"bot/request"
	"bot/response"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/oliveagle/jsonpath"
)

func Mapping(sourceObj string, targetObj string, sourceExpr string, targetExpr string) []byte {

	// Parse JSONPath expressions
	personExpression, err := jsonpath.Compile(sourceExpr)
	if err != nil {
		fmt.Println("Error compiling person JSONPath expression:", err)
		return nil
	}

	dataExpression, err := jsonpath.Compile(targetExpr)
	if err != nil {
		fmt.Println("Error compiling data JSONPath expression:", err)
		return nil
	}

	// Parse JSON data
	var personMap, dataMap map[string]interface{}
	err = json.Unmarshal([]byte(sourceObj), &personMap)
	if err != nil {
		fmt.Printf("Error unmarshaling person JSON: %v\nExpression: %v\nObject: %v", err.Error(), sourceExpr, sourceObj)
		return nil
	}

	err = json.Unmarshal([]byte(targetObj), &dataMap)
	if err != nil {
		fmt.Println("Error unmarshaling data JSON:", err)
		return nil
	}

	// Evaluate JSONPath expression for personMap
	nameResult, err := personExpression.Lookup(personMap)
	if err != nil {
		fmt.Println("Error evaluating person JSONPath expression:", err)
		return nil
	}

	//valida que el campo exista
	oldValue, err := dataExpression.Lookup(dataMap)
	if err != nil {
		fmt.Printf("Error evaluating data JSONPath expression: %v\nExpression: %v\nObject: %v", err.Error(), targetExpr, targetObj)
		return nil
	}
	fmt.Printf("oldValue: %v\n", oldValue)

	dataResult, err := updateField(dataMap, targetExpr, nameResult)
	if err != nil {
		fmt.Println("Error updating data field:", err)
		return nil
	}
	fmt.Printf("dataResult: %v\n", dataResult)

	// Marshal the modified data back to JSON
	modifiedJSON, err := json.Marshal(dataMap)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}

	return modifiedJSON
}

// updateField updates the value of a field in a JSON map using dot-separated field path

func updateField(data interface{}, fieldPath string, value interface{}) (interface{}, error) {
	keys := parseFieldPath(fieldPath)

	// Traverse the data structure
	for i, key := range keys[:len(keys)-1] {
		switch key := key.(type) {
		case string:
			if nested, ok := data.(map[string]interface{}); ok {
				// If it's a map, move to the next level
				data = nested[key]
			} else {
				return nil, fmt.Errorf("nested structure not found for key: %v", key)
			}
		case int:
			// If it's an array, move to the next level
			if arr, ok := data.([]interface{}); ok && key < len(arr) {
				data = arr[key]
			} else {
				return nil, fmt.Errorf("array index out of range: %d", key)
			}
		default:
			return nil, fmt.Errorf("unsupported key type at index %d: %T", i, key)
		}
	}

	// Update the value in the target
	lastKey := keys[len(keys)-1]
	switch lastKey := lastKey.(type) {
	case string:
		if nested, ok := data.(map[string]interface{}); ok {
			// If it's a map, update the value
			nested[lastKey] = value
		} else {
			return nil, fmt.Errorf("nested structure not found for key: %v", lastKey)
		}
	case int:
		// If it's an array, update the value
		if arr, ok := data.([]interface{}); ok && lastKey < len(arr) {
			arr[lastKey] = value
		} else {
			return nil, fmt.Errorf("array index out of range *********** : %d", lastKey)
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %T", lastKey)
	}

	return data, nil
}

func extractArrayKey(input string) (string, error) {
	// Define a regular expression pattern to match the array key between brackets
	re := regexp.MustCompile(`\[(\d+)\]`)

	// Find the first match in the input string
	match := re.FindStringSubmatch(input)

	// Check if there is a match
	if len(match) == 2 {
		// The array key is captured by the first submatch group
		return match[1], nil
	}

	// No match found
	return "", fmt.Errorf("no array key found in the input string")
}
func extractArraName(input string) (string, error) {
	// Expresión regular para extraer la palabra "data"
	re := regexp.MustCompile(`^([^\[\]]+)`)

	// Encontrar coincidencias
	matches := re.FindStringSubmatch(input)

	// Verificar si hay coincidencias
	if len(matches) > 1 {
		word := matches[1]
		// fmt.Println("Palabra extraída:", word)
		return word, nil
	} else {
		fmt.Println("No se encontró ninguna coincidencia.")
	}
	return "", nil
}

// parseFieldPath parses a dot-separated field path into a slice of keys
func parseFieldPath(fieldPath string) []interface{} {
	// delete "$." from "$.data[0].firstName"
	fieldPath = strings.TrimPrefix(fieldPath, "$.")

	var keys []interface{}
	parts := strings.Split(fieldPath, ".")

	for _, part := range parts {
		// fmt.Printf("part: %v\n", part)
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			indexStr, err := extractArrayKey(part)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil
			}
			//extractArraName(part)
			arrayName, err := extractArraName(part)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			keys = append(keys, arrayName)
			keys = append(keys, index)
		} else {
			keys = append(keys, part)
		}
	}

	return keys
}

func ExecMappingByName(mappingName string, requests *request.Requests, responses map[string]response.Response) []byte {
	mappings := requests.Mappings
	output := ""
	for _, mappingData := range mappings {
		if *mappingData.Name == mappingName {
			objects := mappingData.Params.Objects
			sourceObject := responses[*objects.Source].Body.(string)
			output = responses[*objects.Target].Body.(string)
			for _, params := range mappingData.Params.Expressions {
				output = string(Mapping(sourceObject,
					output,
					*params.Source,
					*params.Target))
			}

		}
	}
	return []byte(output)
}
