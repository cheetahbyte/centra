package content

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/helper"
)

func addBinaryFromFile(slug string, path string) error {
	logger := helper.AcquireLogger()
	ext := strings.ToLower(filepath.Ext(path))

	ct := mime.TypeByExtension(ext)

	if ct == "" {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				logger.Error().Err(err).Msg("could close the file correctly")
			}
		}()

		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		ct = http.DetectContentType(buf[:n])
	}

	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = strings.TrimSpace(ct[:i])
	}

	return cache.AddBinaryRef(slug+filepath.Ext(path), ct, path, map[string]any{
		"ext": ext,
	})
}

func handleBinary(key, path string, data []byte) error {
	return addBinaryFromFile(key, path)
}
