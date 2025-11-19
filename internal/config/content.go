package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrNotFound = errors.New("not found")

type ContentRecord map[string]any

func GetContentNode(path string) {}

func GetCollection(collection string) ([]ContentRecord, error) {
	dir := filepath.Join(GetContentRoot(), collection)

	files, err := os.ReadDir(dir)
	if err != nil {
		return []ContentRecord{}, nil
	}

	var entries []ContentRecord

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if !strings.HasSuffix(name, ".yaml") &&
			!strings.HasSuffix(name, ".yml") {
			continue
		}

		full := filepath.Join(dir, name)

		raw, err := os.ReadFile(full)
		if err != nil {
			continue
		}

		var data map[string]any
		if err := yaml.Unmarshal(raw, &data); err != nil {
			continue
		}

		slug := ""
		if v, ok := data["slug"]; ok {
			slug = v.(string)
		} else {
			// remove .yaml or .yml
			slug = strings.TrimSuffix(name, filepath.Ext(name))
		}

		data["slug"] = slug
		entries = append(entries, data)
	}

	return entries, nil
}

func GetEntry(collection, slug string) (ContentRecord, error) {
	dir := filepath.Join(GetContentRoot(), collection)

	possible := []string{
		filepath.Join(dir, slug+".yaml"),
		filepath.Join(dir, slug+".yml"),
	}

	for _, path := range possible {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var data map[string]any
		if err := yaml.Unmarshal(raw, &data); err != nil {
			continue
		}

		data["slug"] = slug
		return data, nil
	}

	return nil, ErrNotFound
}
