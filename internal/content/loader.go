package content

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheetahbyte/centra/internal/cache"
)

// this function iterates over all files and adds them to the store
func LoadAll(contentDir string) error {
	count := 0
	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		if ext == ".yaml" || ext == ".yml" {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			err = cache.AddAndConv(path, b)
			if err != nil {
				return err
			}
			count++
		}

		return nil
	})
	fmt.Println("Indexed", count, "files!")
	return err
}
