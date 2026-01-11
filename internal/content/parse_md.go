package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/goccy/go-yaml"
)

var (
	fmOpen  = []byte("---\n")
	fmClose = []byte("\n---\n")
)

func splitFrontmatter(raw []byte) (meta []byte, body []byte, hasFM bool, err error) {
	if !bytes.HasPrefix(raw, fmOpen) {
		return nil, raw, false, nil
	}

	rest := raw[len(fmOpen):]
	i := bytes.Index(rest, fmClose)
	if i < 0 {
		return nil, nil, false, fmt.Errorf("frontmatter starts with --- but no closing --- found")
	}

	meta = rest[:i]
	body = rest[i+len(fmClose):]
	return meta, body, true, nil
}

func handleMD(key string, path string, data []byte) error {
	return addMarkdown(key, data)
}

func addMarkdown(slug string, raw []byte) error {
	metadata := make(map[string]any)

	metaBytes, bodyBytes, hasFM, err := splitFrontmatter(raw)
	if err != nil {
		return err
	}

	if hasFM {
		if err := yaml.Unmarshal(bytes.TrimSpace(metaBytes), &metadata); err != nil {
			return err
		}
	}

	body := strings.TrimLeft(string(bodyBytes), "\n")

	processedBody := ProcessVariables(body)
	processedMetadata := ProcessMap(metadata)

	doc := map[string]any{
		"body": processedBody,
	}
	for k, v := range processedMetadata {
		doc[k] = v
	}

	jsonData, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	processedMetadata["kind"] = "markdown"
	cache.Insert(slug, processedMetadata, jsonData, "application/json")
	return nil
}
