package utility

import (
	"encoding/json"
)

// minifyJSON takes a string with JSON content and returns a minified version.
func MinifyJSON(jsonStr string) (string, error) {
	var obj interface{}

	// Unmarshal the JSON string into a generic interface.
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return "", err
	}

	// Marshal the object back into a JSON string with no extra whitespace (minified).
	minifiedJSON, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	// Return the minified JSON as a string.
	return string(minifiedJSON), nil
}
