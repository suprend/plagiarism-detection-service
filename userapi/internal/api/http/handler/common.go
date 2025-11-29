package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}

func extractWorkID(path, suffix string) (string, bool) {
	if !strings.HasPrefix(path, "/works/") || !strings.HasSuffix(path, suffix) {
		return "", false
	}
	work := strings.TrimPrefix(path, "/works/")
	work = strings.TrimSuffix(work, suffix)
	work = strings.TrimSuffix(work, "/")
	if work == "" {
		return "", false
	}
	return work, true
}
