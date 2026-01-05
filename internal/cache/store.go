package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"sync"

	"github.com/goccy/go-yaml"
)

var mu sync.RWMutex

func GetNode(path string) *Node {
	mu.RLock()
	defer mu.RLock()

	return ROOT_NODE.Lookup(path)
}

func Get(path string) []byte {
	mu.RLock()
	defer mu.RUnlock()

	node := ROOT_NODE.Lookup(path)
	if node == nil {
		return nil
	}
	return node.GetData()
}

func AddYAML(slug string, raw []byte) error {
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

	mu.Lock()
	defer mu.Unlock()
	ROOT_NODE.Insert(slug, metadata, jsonData)
	return nil
}

func InvalidateAll() {
	mu.Lock()
	defer mu.Unlock()

	ROOT_NODE = NewNode("root")
}

func GetCacheStats() (int, int64) {
	mu.RLock()
	defer mu.RUnlock()

	return ROOT_NODE.calculateStats()
}
