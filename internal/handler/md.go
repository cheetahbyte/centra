package handler

import "github.com/cheetahbyte/centra/internal/ingest"

// markdown handling for files
func handleMD(key string, path string, data []byte) error {
	return ingest.AddMarkdown(key, data)
}
