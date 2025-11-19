package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/content"
	"github.com/cheetahbyte/centra/internal/domain"
	gitadapter "github.com/cheetahbyte/centra/internal/git-adapter"
	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeBinaryJSON(w http.ResponseWriter, status int, v []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(v)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func handleContent(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	path = strings.Trim(path, "/")

	if path == "" {
		writeJSON(w, http.StatusBadRequest, "Invalid content path")
		return
	}
	parts := strings.SplitN(path, "/", 2)
	collection := parts[0]

	if len(parts) == 1 {
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
	} else {
		slug := parts[1]
		if config.GetExperimental("caching") {
			bytes := content.GetEntry(collection, slug)
			writeBinaryJSON(w, 200, bytes)
			return
		}

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
}

func handleWebHook(w http.ResponseWriter, r *http.Request) {
	githubEvent := r.Header.Get("X-Github-Event")
	if githubEvent != "push" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	var whData domain.WebhookData

	if err := json.NewDecoder(r.Body).Decode(&whData); err != nil {
		writeJSON(w, 500, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if whData.Ref != "refs/heads/main" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	contentRoot := config.GetContentRoot()

	w.WriteHeader(http.StatusAccepted)

	go func(wh domain.WebhookData, root string) {
		if err := gitadapter.UpdateRepo(root); err != nil {
			fmt.Println("error updating repo:", err)
			return
		}

		cache.InvalidateAll()

		if err := content.LoadAll(root); err != nil {
			fmt.Println("error loading content:", err)
			return
		}
	}(whData, contentRoot)

}
