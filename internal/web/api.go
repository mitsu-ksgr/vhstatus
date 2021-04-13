package web

import (
	"encoding/json"
	"net/http"
)

func ApiGetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getVHStatusParams())
}
