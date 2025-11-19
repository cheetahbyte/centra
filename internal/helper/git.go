package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	gitadapter "github.com/cheetahbyte/centra/internal/git-adapter"
)

// this function ensures that the given url is a ssh github url
func MakeSSHRepo(url string) string {
	sshBaseUrl := "git@github.com"
	if !strings.HasPrefix(url, "http") {
		if !strings.HasPrefix(url, "ssh") {
			log.Panic("dont know handle this url")
		}
		return url
	}

	fmt.Println("git url is http based. converting.")

	splitUrl := strings.Split(url, "/")

	newUrl := fmt.Sprintf("%s:%s", sshBaseUrl, strings.Join(splitUrl[len(splitUrl)-2:], "/"))

	return newUrl
}

// this function checks if the git repo exists in the directory, that is set as the CONTENT_ROOT
func EnsureRepo(gitRepoUrl, contentDir string) error {
	dir := filepath.Join(contentDir, ".git")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	return gitadapter.CloneRepo(gitRepoUrl, contentDir)
}
