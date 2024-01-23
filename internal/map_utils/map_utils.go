package map_utils

import "net/http"

func MapsEqual(a, b map[string]string) bool {
	for keyA, valueA := range a {
		// Check if the key exists in map b
		if valueB, ok := b[keyA]; ok {
			// If the values are not equal, maps are not equal
			if valueA != valueB {
				return false
			}
		} else {
			// If the key doesn't exist in map b, maps are not equal
			return false
		}
	}
	// All keys in map a exist in map b with equal values
	return true
}

func HeaderToMap(header http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range header {
		// Use the first value, assuming you only want a single string value for each key
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
