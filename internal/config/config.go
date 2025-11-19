package config

import "os"

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
