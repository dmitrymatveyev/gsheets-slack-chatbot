package utility

import (
	"encoding/json"
	"net/http"
)

// WriteInternalError function writes 500 response
func WriteInternalError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(err)
}

// WriteBadRequest function writes 400 response
func WriteBadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}

// WriteResponse function writes 200 response
func WriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
