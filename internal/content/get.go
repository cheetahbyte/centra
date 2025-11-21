package content

import (
	"github.com/cheetahbyte/centra/internal/cache"
)

func Get(path string) *cache.Node {
	return cache.GetNode(path)
}
