package handlers

import (
	"encoding/json"
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	status := map[string]bool{
		"status": true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}
