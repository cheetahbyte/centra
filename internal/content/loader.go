package content

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/logger"
)

func LoadAll(contentDir string) error {
	root := filepath.Clean(contentDir)
	count := 0

	logger := logger.AcquireLogger()

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

	logger.Info().Int("files", count).Msg("cached files into tree")
	nodes, size := cache.GetCacheStats()
	logger.Debug().Int("nodes", nodes).Float64("MB", float64(size)/1024/1024).Msg("tree stats")
	return err
}
