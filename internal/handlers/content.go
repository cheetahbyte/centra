package handlers

import (
	"net/http"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/go-chi/chi/v5"
)

type CollectionItem struct {
	Slug string         `json:"slug"`
	Meta map[string]any `json:"meta"`
}

func HandleContent(w http.ResponseWriter, r *http.Request) {
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
		q := helper.ParseQueryParams(r)
		for p, child := range items {
			meta := child.GetMetadata()
			// this is probably not ideal but it works for now.
			if !helper.MatchesQuery(meta, q) {
				continue
			}

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
