package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetGitRepo() string {
	if root := os.Getenv("GITHUB_REPO_URL"); root != "" {
		return root
	}
	return ""
}

func GetKeysDir() string {
	if dir := os.Getenv("KEYS_DIR"); dir != "" {
		return dir
	}
	return "/keys"
}

func GetPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "3000"
}

func GetContentRoot() string {
	if root := os.Getenv("CONTENT_ROOT"); root != "" {
		return root
	}
	return "/content"
}

func GetCorsAllowedOrigins() []string {
	if raw := os.Getenv("CORS_ALLOWED_ORIGINS"); raw != "" {
		raw = strings.ReplaceAll(raw, "'", `"`)
		var items []string
		if err := json.Unmarshal([]byte(raw), &items); err != nil {
			panic(err)
		}
	}
	return []string{"*"}
}

func GetCorsAllowedMethods() []string {
	if raw := os.Getenv("CORS_ALLOWED_METHODS"); raw != "" {
		raw = strings.ReplaceAll(raw, "'", `"`)
		var items []string
		if err := json.Unmarshal([]byte(raw), &items); err != nil {
			panic(err)
		}
	}
	return []string{"GET", "HEAD", "OPTIONS"}
}

func GetCorsAllowedHeaders() []string {
	if raw := os.Getenv("CORS_ALLOWED_HEADERS"); raw != "" {
		raw = strings.ReplaceAll(raw, "'", `"`)
		var items []string
		if err := json.Unmarshal([]byte(raw), &items); err != nil {
			panic(err)
		}
	}
	return []string{"*"}
}

func GetCorsExposedHeaders() []string {
	if raw := os.Getenv("CORS_EXPOSED_HEADERS"); raw != "" {
		raw = strings.ReplaceAll(raw, "'", `"`)
		var items []string
		if err := json.Unmarshal([]byte(raw), &items); err != nil {
			panic(err)
		}
	}
	return []string{"Cache-Control", "Content-Language", "Content-Length", "Content-Type", "Expires", "Last-Modified"}
}

func GetCorsMaxAge() int {
	if raw := os.Getenv("CORS_MAX_AGE"); raw != "" {
		numb, err := strconv.Atoi(raw)
		if err != nil {
			panic(err)
		}
		return numb
	}
	return 360
}

func GetCorsAllowCredentials() bool {
	raw := os.Getenv("CORS_ALLOW_CREDENTIALS")
	switch raw {
	case "true":
		return true
	default:
		return false
	}
}

func GetPrivateSSHKey() string {
	raw := os.Getenv("SSH_PRIVATE_KEY")
	return raw
}

func GetPublicSSHKey() string {
	raw := os.Getenv("SSH_PUBLIC_KEY")
	return raw
}

// this generic function returns the raw object
func GetExperimental(featureName string) bool {
	raw := os.Getenv(fmt.Sprintf("EXPERIMENTAL_%s", strings.ToUpper(featureName)))
	switch raw {
	case "true":
		return true
	default:
		return false
	}
}
