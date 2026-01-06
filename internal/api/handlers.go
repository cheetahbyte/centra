package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/content"
	"github.com/cheetahbyte/centra/internal/domain"
	gitadapter "github.com/cheetahbyte/centra/internal/git-adapter"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/cheetahbyte/centra/internal/logger"
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

type CollectionItem struct {
	Slug string         `json:"slug"`
	Meta map[string]any `json:"meta"`
}

func handleContent(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	path = strings.Trim(path, "/")

	if path == "" {
		writeJSON(w, http.StatusBadRequest, "Invalid content path")
		return
	}

	node := cache.GetNode(path)
	if node == nil {
		writeJSON(w, http.StatusNotFound, "Content not found")
		return
	}

	if !node.IsLeaf() {
		items := node.GetChildren()
		collectionItems := make([]CollectionItem, 0, len(items))
		for p, child := range items {
			collectionItems = append(collectionItems, CollectionItem{
				Slug: path + "/" + p,
				Meta: child.GetMetadata(),
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"collection": path,
			"items":      collectionItems,
		})
		return
	}

	meta := node.GetMetadata()
	ct, _ := meta["contentType"].(string)
	if ct == "" {
		ct = "application/octet-stream"
	}

	// check if its a binary file
	if fp := node.GetFilePath(); fp != "" {
		w.Header().Set("Content-Type", ct)
		// maybe caching

		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.ServeFile(w, r, fp)
		return
	}

	// serve memory bytes
	data := node.GetData()
	if data == nil {
		writeJSON(w, http.StatusNotFound, "Content not found")
		return
	}

	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func handleWebHook(w http.ResponseWriter, r *http.Request) {
	log := logger.AcquireLogger()

	signatureHeader := r.Header.Get("X-Hub-Signature-256")

	githubEvent := r.Header.Get("X-Github-Event")
	if githubEvent != "push" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("could not read request body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	webhookSecret := config.GetWebhookSecret()
	if webhookSecret != "" {
		if err := helper.VerifySignature(bodyBytes, signatureHeader, webhookSecret); err != nil {
			log.Error().Err(err).Msg("invalid webhook signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var whData domain.WebhookData
	if err := json.Unmarshal(bodyBytes, &whData); err != nil {
		log.Error().Err(err).Msg("invalid json body")
		writeJSON(w, 400, map[string]any{
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
			log.Error().Err(err).Msg("error during repo update")
			return
		}

		cache.InvalidateAll()

		if err := content.LoadAll(root); err != nil {
			log.Error().Err(err).Msg("failed to cache during webhook update")
			return
		}
	}(whData, contentRoot)
}
