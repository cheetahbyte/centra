package cache

import (
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

func AddAndConv(slug string, yamlData []byte) error {
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}

	Add(slug, jsonData)
	return nil
}

func Add(slug string, jsonData []byte) error {
	mu.Lock()
	defer mu.Unlock()

	ROOT_NODE.Insert(slug, jsonData)
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
