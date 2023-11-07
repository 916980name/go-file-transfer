package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func HttpReadBody(r *http.Request, customType interface{}) error {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return err
	}

	// Create an instance of the custom type

	// Unmarshal the request body into the custom type
	err = json.Unmarshal(body, &customType)
	if err != nil {
		// http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return err
	}
	return nil
}
