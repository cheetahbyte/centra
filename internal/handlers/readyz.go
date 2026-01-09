package handlers

import (
	"net/http"

	"github.com/cheetahbyte/centra/internal/cache"
)

func HandleReadyz(w http.ResponseWriter, r *http.Request) {
	if !cache.IsReady() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "initializing",
		})
		return
	}

	cachedNodes, treeSize := cache.ROOT_NODE.CalculateStats()
	writeJSON(w, http.StatusOK, map[string]any{
		"status":      "ok",
		"cachedNodes": cachedNodes,
		"treeSize":    treeSize,
	})
}
