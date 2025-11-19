package content

import (
	"github.com/cheetahbyte/centra/internal/cache"
)

func GetEntry(collection, slug string) []byte {
	lookupPath := collection + "/" + slug
	bytes := cache.Get(lookupPath)
	return bytes
}
