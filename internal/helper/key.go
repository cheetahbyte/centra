package helper

import (
	"os"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/logger"
)

func PrettyKey(pubKey string) {
	logger := logger.AcquireLogger()
	if config.GetPublicSSHKey() == "" {
		b, _ := os.ReadFile(pubKey)
		logger.Info().Str("public key", string(b)).Msg("Add the deploy key to your github repository.")
	}
}
