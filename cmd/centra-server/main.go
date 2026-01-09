package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/api"
	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/content"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/cheetahbyte/drift/keys"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	api.Register(r)

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log := helper.AcquireLogger()

	if conf.GitRepo != "" {
		if conf.PublicKey != "" {
			log.Info().Str("ssh key", conf.PublicKey).Msg("add this ssh key to your github repository as deploy key.")
		} else {
			publicKeyPath := filepath.Join(conf.KeysDir, "id_ed25519.pub")
			pubKey, err := keys.GetPublicKey(publicKeyPath)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to output public ssh key")
			}
			log.Info().Str("ssh key", helper.PublicKeyToString(pubKey)).Msg("add this ssh key to your github repository as deploy key.")
		}
		gitClient := helper.SetupGit()
		if err := gitClient.Prepare(conf.GitRepo, conf.ContentRoot); err != nil {
			log.Fatal().Err(err).Msg("failed to clone or prepare git repository")
		}

	}

	go func() {
		if err := content.LoadAll(conf.ContentRoot); err != nil {
			log.Fatal().Err(err).Msg("caching did not work.")
		}
		cache.SetReady(true)
		log.Info().Msg("content loaded. service is now ready.")
	}()

	log.Info().Str("port", conf.Port).Msg("centra api is running.")

	err = http.ListenAndServe(":"+conf.Port, r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server.")
	}
}
