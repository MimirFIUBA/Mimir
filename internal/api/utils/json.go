package utils

import (
	"encoding/json"
	"mimir/internal/api/models"
	"net/http"
)

// TODO(#19) - Improve error handling
// TODO - Handle nil cases
// TODO - Handle empty cases
// TODO - Handle [] cases
func RespondWithJSONItems(w http.ResponseWriter, code int, payload interface{}) error {
	itemsResponse := models.ItemsResponse{
		Status: code,
		Items:  []any{payload},
	}
	response, err := json.Marshal(itemsResponse)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}
