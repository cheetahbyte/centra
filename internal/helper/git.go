package helper

import (
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/logger"
	"github.com/cheetahbyte/drift/git"
	"github.com/cheetahbyte/drift/keys"
)

func SetupGit() *git.Client {
	log := logger.AcquireLogger()
	keysDir := config.GetKeysDir()

	pubKeyPath, err := keys.Setup(
		keysDir,
		config.GetPrivateSSHKey(),
		config.GetPublicSSHKey(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup ssh keys")
	}

	if config.GetPublicSSHKey() == "" {
		log.Info().Str("path", pubKeyPath).Msg("SSH public key ready")
	}

	privateKeyPath := filepath.Join(keysDir, "id_ed25519")
	return git.New(privateKeyPath)
}
