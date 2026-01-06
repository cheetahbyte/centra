package main

import (
	"log"
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

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.AcquireLogger()

	if conf.GitRepo != "" {
		gitClient := helper.SetupGit()
		if err := gitClient.Prepare(conf.GitRepo, conf.ContentRoot); err != nil {
			log.Fatal().Err(err).Msg("failed to clone or prepare git repository")
		}
	}

	if err := content.LoadAll(conf.ContentRoot); err != nil {
		log.Fatal().Err(err).Msg("caching did not work.")
	}

	log.Info().Str("port", conf.Port).Msg("centra api is running.")

	err = http.ListenAndServe(":"+conf.Port, r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server.")
	}
}
