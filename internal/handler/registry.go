package handler

import (
	"strings"
)

type FileHandler func(key string, path string, data []byte) error

var handlers = map[string]FileHandler{
	".yaml": handleYaml,
	".yml":  handleYaml,
	".md":   handleMD,
}

var BinaryAllowList = map[string]bool{
	// images
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true, ".svg": true, ".ico": true, ".avif": true,
	// docs
	".pdf": true,
	".mp4": true, ".webm": true, ".mp3": true, ".wav": true, ".ogg": true,
	".woff": true, ".woff2": true, ".ttf": true, ".otf": true,
	".zip": true,
}

func HandleFor(ext string) FileHandler {
	ext = strings.ToLower(ext)

	if h, ok := handlers[ext]; ok {
		return h
	}
	if BinaryAllowList[ext] {
		return handleBinary
	}
	return handleIgnore // unknown/unwanted types
}

func handleIgnore(key, path string, data []byte) error { return nil }
