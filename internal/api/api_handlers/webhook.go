package apihandlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/content"
	"github.com/cheetahbyte/centra/internal/logger"
	"github.com/cheetahbyte/drift/git"
	"github.com/cheetahbyte/drift/keys"
	"github.com/cheetahbyte/drift/webhook"
)

func setupGit() *git.Client {
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

func HandleWebHook(w http.ResponseWriter, r *http.Request) {
	log := logger.AcquireLogger()

	if r.Header.Get("X-Github-Event") != "push" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("could not read webhook request body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	webhookSecret := config.GetWebhookSecret()
	signatureHeader := r.Header.Get("X-Hub-Signature-256")

	if webhookSecret != "" {
		if err := webhook.VerifySignature(bodyBytes, signatureHeader, webhookSecret); err != nil {
			log.Error().Err(err).Msg("unauthorized: invalid webhook signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var payload struct {
		Ref string `json:"ref"`
	}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Error().Err(err).Msg("invalid json body in webhook")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if payload.Ref != "refs/heads/main" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	contentRoot := config.GetContentRoot()

	go func(root string) {
		gitClient := setupGit()
		if gitClient == nil {
			log.Error().Err(err).Msg("aborting update: failed to setup git client")
			return
		}

		// Pull changes
		changedFiles, err := gitClient.Pull(root, "main")
		if err != nil {
			log.Error().Err(err).Msg("failed to update repository")
			return
		}

		log.Info().Int("files_changed", changedFiles).Msg("git pull successful")

		cache.InvalidateAll()
		if err := content.LoadAll(root); err != nil {
			log.Error().Err(err).Msg("failed to reload content after git update")
			return
		}

		log.Info().Msg("content update complete")

	}(contentRoot)
}
