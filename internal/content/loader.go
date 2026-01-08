package content

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/handler"
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
			switch d.Name() {
			case ".git", "node_modules", ".next", "dist", "build":
				return fs.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		h := handler.HandleFor(ext)

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(root, strings.TrimSuffix(path, ext))
		key := filepath.ToSlash(rel)

		if err := h(key, path, b); err != nil {
			logger.Error().Err(err).Str("path", path).Msg("failed to cache file")
			return nil
		}
		logger.Debug().Str("path", path).Msg("visited file")

		count++
		return nil
	})

	logger.Info().Int("files", count).Msg("cached files into tree")
	nodes, size := cache.GetCacheStats()
	logger.Debug().Int("nodes", nodes).Float64("MB", float64(size)/1024/1024).Msg("tree stats")
	return err
}
