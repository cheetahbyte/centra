package content

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/logger"
)

// this function iterates over all files and adds them to the store
func LoadAll(contentDir string) error {
	root := filepath.Clean(contentDir)

	count := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		base := strings.TrimSuffix(path, ext)

		// get path relative to content dir: e.g. "pages/home"
		rel, err := filepath.Rel(root, base)
		if err != nil {
			return err
		}

		key := filepath.ToSlash(rel)

		if err := cache.AddAndConv(key, b); err != nil {
			return err
		}

		count++
		return nil
	})

	logger := logger.AcquireLogger()
	logger.Info().Int("files", count).Msg("cached files")
	return err
}
