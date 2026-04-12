package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func UnmarshalJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid json body", fmt.Errorf("failed to unmarshal body: %w", err))
		return false
	}
	return true
}

func WriteErrorJSON(w http.ResponseWriter, status int, clientErr string, logErr error) {
	WriteJSON(w, status, map[string]string{
		"client_err": clientErr,
		"log_err":    logErr.Error(),
	})
}

func WriteNotFoundErrorJSON(w http.ResponseWriter, logErr error) {
	WriteErrorJSON(
		w, http.StatusNotFound,
		"there is no ressource assosiated with provided id",
		logErr,
	)
}

func WriteInternalServerError(w http.ResponseWriter, logErr error) {
	WriteErrorJSON(w, http.StatusInternalServerError, "internal server error", logErr)
}
