package ingest

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
)

func AddBinaryFromFile(slug string, path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	ct := mime.TypeByExtension(ext)

	if ct == "" {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		ct = http.DetectContentType(buf[:n])
	}

	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = strings.TrimSpace(ct[:i])
	}

	return cache.AddBinaryRef(slug, ct, path, map[string]any{
		"ext": ext,
	})
}
