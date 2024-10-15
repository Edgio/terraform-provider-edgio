package utility

import (
	"encoding/json"
)

func MinifyJSON(jsonStr string) (string, error) {
	var obj interface{}

	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return "", err
	}

	// Marshal the object back into a JSON string with no extra whitespace (minified).
	minifiedJSON, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(minifiedJSON), nil
}
