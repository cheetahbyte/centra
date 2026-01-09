package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/goccy/go-yaml"
)

func handleYaml(key string, path string, data []byte) error {
	return addYaml(key, data)
}

func addYaml(slug string, raw []byte) error {
	metadata := make(map[string]any)
	bodyMap := make(map[string]any)
	fullData := make(map[string]any)

	parts := bytes.SplitN(raw, []byte("---\n"), 2)

	switch len(parts) {
	case 2:
		if err := yaml.Unmarshal(bytes.TrimSpace(parts[0]), &metadata); err != nil {
			return err
		}
		if err := yaml.Unmarshal(bytes.TrimSpace(parts[1]), &bodyMap); err != nil {
			return err
		}
		maps.Copy(fullData, metadata)
		maps.Copy(fullData, bodyMap)

	case 1:
		if err := yaml.Unmarshal(bytes.TrimSpace(parts[0]), &bodyMap); err != nil {
			return err
		}
		maps.Copy(fullData, bodyMap)

	default:
		return fmt.Errorf("invalid yaml format")
	}

	jsonData, err := json.Marshal(fullData)
	if err != nil {
		return err
	}

	metadata["kind"] = "yaml"
	cache.Insert(slug, metadata, jsonData, "application/json")
	return nil
}
