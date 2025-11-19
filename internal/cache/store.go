package cache

import (
	"sync"

	"github.com/goccy/go-yaml"
)

type ContentStore struct {
	mu        sync.RWMutex
	jsonBytes map[string][]byte
}

var store = &ContentStore{
	jsonBytes: make(map[string][]byte),
}

func Get(slug string) []byte {
	store.mu.RLock()
	data := store.jsonBytes[slug]
	store.mu.RUnlock()
	return data
}

func AddAndConv(slug string, yamlData []byte) error {
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}

	if err := Add(slug, jsonData); err != nil {
		return err
	}

	return nil
}

func Add(slug string, jsonData []byte) error {
	store.mu.Lock()
	store.jsonBytes[slug] = jsonData
	store.mu.Unlock()

	return nil
}

func InvalidateAll() {
	store.mu.Lock()
	store.jsonBytes = nil
	store.jsonBytes = make(map[string][]byte)
	store.mu.Unlock()
}
