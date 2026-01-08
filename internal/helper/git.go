package helper

import (
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/logger"
	"github.com/cheetahbyte/drift/git"
	"github.com/cheetahbyte/drift/keys"
)

func SetupGit() *git.Client {
	conf := config.Get()
	log := logger.AcquireLogger()

	pubKeyPath, err := keys.Setup(
		conf.KeysDir,
		conf.PrivateKey,
		conf.PublicKey,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup ssh keys")
	}

	if conf.PublicKey == "" {
		log.Info().Str("path", pubKeyPath).Msg("SSH public key ready")
	}

	privateKeyPath := filepath.Join(conf.KeysDir, "id_ed25519")
	return git.New(privateKeyPath)
}
