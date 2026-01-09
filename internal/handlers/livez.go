package handlers

import "net/http"

func HandleLivez(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"status": "ok",
	})
}
