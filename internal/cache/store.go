package cache

import (
	"os"
	"sync"
)

var mu sync.RWMutex

func GetNode(path string) *Node {
	mu.RLock()
	defer mu.RUnlock() // FIX: was mu.RLock()

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

// Optional helper if you want metadata too
func GetWithMetadata(path string) (*Node, []byte, map[string]any) {
	mu.RLock()
	defer mu.RUnlock()

	node := ROOT_NODE.Lookup(path)
	if node == nil {
		return nil, nil, nil
	}
	return node, node.GetData(), node.GetMetadata()
}

func AddBinaryRef(slug string, contentType string, absPath string, metadata map[string]any) error {
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	if metadata == nil {
		metadata = map[string]any{}
	}
	metadata["kind"] = "binary_ref"
	metadata["contentType"] = contentType
	metadata["size"] = fi.Size()
	metadata["mtime"] = fi.ModTime().Unix()

	mu.Lock()
	defer mu.Unlock()

	ROOT_NODE.Insert(slug, metadata, nil, contentType)
	if n := ROOT_NODE.Lookup(slug); n != nil {
		n.typ = contentType
		n.filePath = absPath
	}
	return nil
}

func Insert(slug string, metadata map[string]any, data []byte, typ string) {
	if metadata == nil {
		metadata = map[string]any{}
	}

	metadata["contentType"] = typ
	metadata["size"] = len(data)

	mu.Lock()
	defer mu.Unlock()

	ROOT_NODE.Insert(slug, metadata, data, typ)
	if n := ROOT_NODE.Lookup(slug); n != nil {
		n.typ = typ
	}
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
