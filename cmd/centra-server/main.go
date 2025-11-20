package main

import (
	"net/http"

	"github.com/cheetahbyte/centra/internal/api"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/content"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/cheetahbyte/centra/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	api.Register(r)

	port := config.GetPort()

	logger := logger.AcquireLogger()

	keyDir := config.GetKeysDir()

	pubKey, err := helper.EnsureKeys(keyDir)
	if err != nil {
		logger.Fatal().Err(err).Msg("problem with ssh keys")
	}

	helper.PrettyKey(pubKey)

	repo := config.GetGitRepo()
	if repo != "" {
		helper.EnsureRepo(repo, config.GetContentRoot())
	}

	if config.GetExperimental("caching") {
		if err := content.LoadAll(config.GetContentRoot()); err != nil {
			logger.Fatal().Err(err).Msg("caching did not work.")
		}
	}

	logger.Info().Str("port", port).Msg("centra api is running.")

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server.")
	}
}
