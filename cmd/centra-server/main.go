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
	log := logger.AcquireLogger()

	repoURL := config.GetGitRepo()

	if repoURL != "" {
		gitClient := helper.SetupGit()

		contentRoot := config.GetContentRoot()
		if err := gitClient.Prepare(repoURL, contentRoot); err != nil {
			log.Fatal().Err(err).Msg("failed to clone or prepare git repository")
		}
	}

	if err := content.LoadAll(config.GetContentRoot()); err != nil {
		log.Fatal().Err(err).Msg("caching did not work.")
	}

	log.Info().Str("port", port).Msg("centra api is running.")

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server.")
	}
}
