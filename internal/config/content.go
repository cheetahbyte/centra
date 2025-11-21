package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidPath = errors.New("invalid path")
)

type ContentRecord map[string]any

func safeJoin(root string, parts ...string) (string, error) {
	cleanRoot := filepath.Clean(root)
	sep := string(os.PathSeparator)

	safeParts := []string{cleanRoot}

	for _, p := range parts {
		if p == "" {
			continue
		}

		if filepath.IsAbs(p) {
			return "", ErrInvalidPath
		}

		p = filepath.Clean(p)

		if p == ".." || strings.HasPrefix(p, ".."+sep) {
			return "", ErrInvalidPath
		}

		safeParts = append(safeParts, p)
	}

	joined := filepath.Join(safeParts...)

	rootWithSep := cleanRoot
	if !strings.HasSuffix(rootWithSep, sep) {
		rootWithSep += sep
	}

	if joined != cleanRoot && !strings.HasPrefix(joined+sep, rootWithSep) {
		return "", ErrInvalidPath
	}

	return joined, nil
}

func GetContentNode(path string) (ContentRecord, error) {
	full, err := safeJoin(GetContentRoot(), path)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var data map[string]any
	if err := yaml.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func GetCollection(collection string) ([]ContentRecord, error) {
	dir, err := safeJoin(GetContentRoot(), collection)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []ContentRecord{}, nil
		}
		return nil, err
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

		full, err := safeJoin(dir, name)
		if err != nil {
			continue
		}

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
			if s, ok := v.(string); ok {
				slug = s
			}
		}
		if slug == "" {
			slug = strings.TrimSuffix(name, filepath.Ext(name))
		}

		data["slug"] = slug
		entries = append(entries, data)
	}

	return entries, nil
}

func GetEntry(collection, slug string) (ContentRecord, error) {
	dir, err := safeJoin(GetContentRoot(), collection)
	if err != nil {
		return nil, err
	}

	candidates := []string{
		slug + ".yaml",
		slug + ".yml",
	}

	for _, c := range candidates {
		path, err := safeJoin(dir, c)
		if err != nil {
			continue
		}

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
