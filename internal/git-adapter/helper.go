package gitadapter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/config"
)

func UpdateRepo(dir string) error {
	pullOrigin := "origin"
	pullBranch := "main"

	privateKey := filepath.Join(config.GetKeysDir(), "id_ed25519")

	cmd := exec.Command("git", "pull", pullOrigin, pullBranch, "--ff-only")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	env := os.Environ()
	env = append(env,
		"GIT_SSH_COMMAND=ssh -i "+privateKey+" -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new",
	)
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %v\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}

	fmt.Println("git pull ok: ", stdout.String())
	return nil
}

func CloneRepo(repoUrl string, dist string) error {
	privateKey := filepath.Join(config.GetKeysDir(), "id_ed25519")

	cmd := exec.Command("git", "clone", repoUrl, dist)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	env := os.Environ()
	env = append(env, "GIT_SSH_COMMAND=ssh -i "+privateKey+" -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new")
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %v\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}

	return nil
}
