package handler

import (
	"github.com/cheetahbyte/centra/internal/ingest"
)

// handle for everything that is not yaml or md
func handleBinary(key, path string, data []byte) error {
	return ingest.AddBinaryFromFile(key, path)
}
