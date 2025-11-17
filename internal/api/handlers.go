package api

import (
	"encoding/json"
	"net/http"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func handleGetCollection(w http.ResponseWriter, r *http.Request) {
	collection := chi.URLParam(r, "collection")
	items, err := config.GetCollection(collection)
	if err != nil {
		writeJSON(w, 500, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]any{
		"collection": collection,
		"items":      items,
	})
}

func handleGetEntry(w http.ResponseWriter, r *http.Request) {
	collection := chi.URLParam(r, "collection")
	slug := chi.URLParam(r, "slug")

	entry, err := config.GetEntry(collection, slug)
	if err != nil {
		if err == config.ErrNotFound {
			writeJSON(w, 404, map[string]any{
				"error":      "Not found",
				"collection": collection,
				"slug":       slug,
			})
			return
		}
		writeJSON(w, 500, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, 200, entry)
}
