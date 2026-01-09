package content

import (
	"strings"

	"github.com/cheetahbyte/centra/internal/config"
)

type FileHandler func(key string, path string, data []byte) error

var handlers = map[string]FileHandler{
	".yaml": handleYaml,
	".yml":  handleYaml,
	".md":   handleMD,
}

func HandleFor(ext string) FileHandler {
	conf := config.Get()
	ext = strings.ToLower(ext)

	if h, ok := handlers[ext]; ok {
		return h
	}
	if conf.AnyBinaries || config.BinaryAllowList[ext] {
		return handleBinary
	}
	return handleIgnore // unknown/unwanted types
}

func handleIgnore(key, path string, data []byte) error { return nil }
