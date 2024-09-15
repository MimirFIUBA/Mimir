package responses

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func SendJSONResponse(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Used by middlewares
	w.Header().Set("status", strconv.Itoa(code))

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		return err
	}

	return nil
}

func SendErrorResponse(w http.ResponseWriter, code int, error ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Used by middlewares
	w.Header().Set("status", strconv.Itoa(code))

	jsonResponse, err := json.Marshal(error)
	if err != nil {
		// This should never execute
		panic(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		// This should never execute
		panic(err)
	}
}
