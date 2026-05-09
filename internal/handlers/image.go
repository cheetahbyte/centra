package handlers

import (
	"crypto/sha256"
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/disintegration/imaging"
	"golang.org/x/sync/singleflight"
)

var transformGroup singleflight.Group

type TransformParams struct {
	Width   int
	Height  int
	Quality int
	Format  string
}

func isImageType(ct string) bool {
	switch ct {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	}
	return false
}

func parseTransformParams(r *http.Request) (TransformParams, bool) {
	q := r.URL.Query()
	p := TransformParams{Quality: 80}

	p.Width, _ = strconv.Atoi(q.Get("w"))
	p.Height, _ = strconv.Atoi(q.Get("h"))
	if v, err := strconv.Atoi(q.Get("q")); err == nil && v > 0 {
		p.Quality = v
	}
	p.Format = strings.ToLower(q.Get("format"))

	if p.Width == 0 && p.Height == 0 {
		return TransformParams{}, false
	}
	return p, true
}

func TransformImage(w http.ResponseWriter, r *http.Request, originalPath string, originalCT string, params TransformParams) {
	log := helper.AcquireLogger()
	conf := config.Get()

	maxDim := conf.ImageMaxDim
	if params.Width > maxDim {
		params.Width = maxDim
	}
	if params.Height > maxDim {
		params.Height = maxDim
	}

	outFormat, outCT := resolveFormat(params.Format, originalCT)

	cacheKey := fmt.Sprintf("%x", sha256.Sum256(
		fmt.Appendf(nil, "%s|%d|%d|%d|%s", originalPath, params.Width, params.Height, params.Quality, outFormat),
	))
	cachePath := filepath.Join(conf.ImageCacheDir, cacheKey+"."+outFormat)

	if _, err := os.Stat(cachePath); err == nil {
		w.Header().Set("Content-Type", outCT)
		http.ServeFile(w, r, cachePath)
		return
	}

	_, err, _ := transformGroup.Do(cacheKey, func() (any, error) {
		if err := os.MkdirAll(conf.ImageCacheDir, 0755); err != nil {
			return nil, err
		}

		src, err := decodeImage(originalPath, originalCT)
		if err != nil {
			return nil, err
		}

		bounds := src.Bounds()
		w := params.Width
		h := params.Height
		if w == 0 {
			w = bounds.Dx()
		}
		if h == 0 {
			h = bounds.Dy()
		}
		dst := imaging.Fit(src, w, h, imaging.Linear)

		return nil, encodeImage(dst, cachePath, outFormat, params.Quality)
	})

	if err != nil {
		log.Error().Err(err).Str("path", originalPath).Msg("failed to transform image")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", outCT)
	http.ServeFile(w, r, cachePath)
}

func resolveFormat(requested, originalCT string) (ext, ct string) {
	switch requested {
	case "webp":
		return "webp", "image/webp"
	case "png":
		return "png", "image/png"
	case "gif":
		return "gif", "image/gif"
	case "jpeg", "jpg":
		return "jpg", "image/jpeg"
	}
	// fall back to original format
	switch originalCT {
	case "image/png":
		return "png", "image/png"
	case "image/gif":
		return "gif", "image/gif"
	case "image/webp":
		return "webp", "image/webp"
	default:
		return "jpg", "image/jpeg"
	}
}

func decodeImage(path, ct string) (image.Image, error) {
	if ct == "image/webp" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return webp.Decode(f)
	}
	return imaging.Open(path)
}

func encodeImage(img image.Image, destPath, format string, quality int) error {
	tmp := destPath + ".tmp"

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	var encErr error
	switch format {
	case "webp":
		encErr = webp.Encode(f, img, &webp.Options{Lossless: false, Quality: float32(quality)})
	case "png":
		encErr = imaging.Encode(f, img, imaging.PNG)
	case "gif":
		encErr = imaging.Encode(f, img, imaging.GIF)
	default:
		encErr = imaging.Encode(f, img, imaging.JPEG, imaging.JPEGQuality(quality))
	}

	f.Close()
	if encErr != nil {
		_ = os.Remove(tmp)
		return encErr
	}
	return os.Rename(tmp, destPath)
}
