package utils

import "strings"

func GetValueFromJSON(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	var value interface{} = data

	for _, key := range keys {
		// Verificamos si el valor actual es un mapa
		if mapVal, ok := value.(map[string]interface{}); ok {
			value = mapVal[key]
		} else {
			return nil, false
		}
	}

	return value, true
}
