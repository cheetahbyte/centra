package handler

import (
	"github.com/cheetahbyte/centra/internal/ingest"
)

func handleYaml(key string, path string, data []byte) error {
	return ingest.AddYAML(key, data)
}
