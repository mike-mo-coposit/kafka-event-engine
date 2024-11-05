package utils

import (
	"encoding/json"
	"fmt"
)

// ConvertToType converts an event of type `interface{}` into the specified target type.
func ConvertToType(event interface{}, target interface{}) error {
	var jsonData []byte
	var err error

	// Check if event is a string and convert to []byte
	if strData, ok := event.(string); ok {
		jsonData = []byte(strData)
	} else if byteData, ok := event.([]byte); ok { // If []byte, use it directly
		jsonData = byteData
	} else {
		return fmt.Errorf("event is neither string nor []byte")
	}

	// Unmarshal JSON into target
	err = json.Unmarshal(jsonData, target)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return nil
}
