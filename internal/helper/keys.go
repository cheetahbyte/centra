package helper

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/logger"
	xssh "golang.org/x/crypto/ssh"
)

func EnsureKeys(keysDir string) (string, error) {
	log := logger.AcquireLogger()

	if err := os.MkdirAll(keysDir, 0o700); err != nil {
		log.Debug().
			Str("keysDir", keysDir).
			Err(err).
			Msg("failed to create keys directory")
		return "", fmt.Errorf("create keys dir: %w", err)
	}
	log.Debug().
		Str("keysDir", keysDir).
		Msg("ensured keys directory exists")

	privateKeyPath := filepath.Join(keysDir, "id_ed25519")
	publicKeyPath := filepath.Join(keysDir, "id_ed25519.pub")

	configPrivKey := config.GetPrivateSSHKey()
	configPubKey := config.GetPublicSSHKey()
	// Case 1: keys are provided via config
	if configPrivKey != "" {
		log.Debug().
			Str("privateKeyPath", privateKeyPath).
			Str("publicKeyPath", publicKeyPath).
			Msg("using SSH keys from config")

		privKeyBytes := []byte(configPrivKey)
		if err := os.WriteFile(privateKeyPath, privKeyBytes, 0o600); err != nil {
			log.Debug().
				Str("privateKeyPath", privateKeyPath).
				Err(err).
				Msg("failed to write config private key")
			return "", fmt.Errorf("failed to write config private key: %w", err)
		}

		if configPubKey != "" {
			log.Debug().
				Str("publicKeyPath", publicKeyPath).
				Msg("writing public key from config")
			if err := os.WriteFile(publicKeyPath, []byte(configPubKey), 0o644); err != nil {
				log.Debug().
					Str("publicKeyPath", publicKeyPath).
					Err(err).
					Msg("failed to write config public key")
				return "", fmt.Errorf("failed to write config public key: %w", err)
			}
		} else {
			log.Debug().
				Str("publicKeyPath", publicKeyPath).
				Msg("deriving public key from config private key")
			if err := deriveAndSavePublicKey(privKeyBytes, publicKeyPath); err != nil {
				log.Debug().
					Str("publicKeyPath", publicKeyPath).
					Err(err).
					Msg("failed to derive public key from config private key")
				return "", fmt.Errorf("failed to derive public key from config: %w", err)
			}
		}

		log.Debug().
			Str("publicKeyPath", publicKeyPath).
			Msg("SSH keys from config are ready")
		return publicKeyPath, nil
	}

	// Case 2: private key already exists on disk
	if _, err := os.Stat(privateKeyPath); err == nil {
		log.Debug().
			Str("privateKeyPath", privateKeyPath).
			Msg("found existing private key on disk")

		if _, err := os.Stat(publicKeyPath); err == nil {
			log.Debug().
				Str("publicKeyPath", publicKeyPath).
				Msg("found existing public key on disk")
			return publicKeyPath, nil
		}

		log.Debug().
			Str("privateKeyPath", privateKeyPath).
			Str("publicKeyPath", publicKeyPath).
			Msg("public key missing; deriving from existing private key")

		existingPrivBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			log.Debug().
				Str("privateKeyPath", privateKeyPath).
				Err(err).
				Msg("failed to read existing private key")
			return "", fmt.Errorf("read existing private key: %w", err)
		}
		if err := deriveAndSavePublicKey(existingPrivBytes, publicKeyPath); err != nil {
			log.Debug().
				Str("publicKeyPath", publicKeyPath).
				Err(err).
				Msg("failed to derive public key from existing private key")
			return "", fmt.Errorf("derive public key from existing file: %w", err)
		}

		log.Debug().
			Str("publicKeyPath", publicKeyPath).
			Msg("derived public key from existing private key")
		return publicKeyPath, nil
	}

	// Case 3: no keys -> generate new pair
	log.Debug().
		Str("privateKeyPath", privateKeyPath).
		Str("publicKeyPath", publicKeyPath).
		Msg("no existing keys found; generating new ed25519 keypair")

	return generateAndSaveKeys(privateKeyPath, publicKeyPath)
}

func deriveAndSavePublicKey(privKeyBytes []byte, pubPath string) error {
	log := logger.AcquireLogger()

	log.Debug().
		Str("publicKeyPath", pubPath).
		Msg("parsing private key to derive public key")

	signer, err := xssh.ParsePrivateKey(privKeyBytes)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to parse private key while deriving public key")
		return fmt.Errorf("parse private key: %w", err)
	}

	sshPubKey := signer.PublicKey()
	pubBytes := xssh.MarshalAuthorizedKey(sshPubKey)

	if err := os.WriteFile(pubPath, pubBytes, 0o644); err != nil {
		log.Debug().
			Str("publicKeyPath", pubPath).
			Err(err).
			Msg("failed to write derived public key")
		return fmt.Errorf("write derived public key: %w", err)
	}

	log.Debug().
		Str("publicKeyPath", pubPath).
		Msg("successfully derived and wrote public key")

	return nil
}

func generateAndSaveKeys(privPath, pubPath string) (string, error) {
	log := logger.AcquireLogger()

	log.Debug().
		Str("privateKeyPath", privPath).
		Str("publicKeyPath", pubPath).
		Msg("generating new ed25519 keypair")

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to generate ed25519 keypair")
		return "", fmt.Errorf("generate ed25519 key: %w", err)
	}

	privBlock, err := xssh.MarshalPrivateKey(priv, "")
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to marshal private key")
		return "", fmt.Errorf("marshal private key: %w", err)
	}

	privPEMBytes := pem.EncodeToMemory(privBlock)

	if err := os.WriteFile(privPath, privPEMBytes, 0o600); err != nil {
		log.Debug().
			Str("privateKeyPath", privPath).
			Err(err).
			Msg("failed to write private key")
		return "", fmt.Errorf("write private key: %w", err)
	}

	log.Debug().
		Str("privateKeyPath", privPath).
		Msg("wrote private key to disk")

	sshPubKey, err := xssh.NewPublicKey(pub)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to create public key from generated keypair")
		return "", fmt.Errorf("create public key: %w", err)
	}

	pubBytes := xssh.MarshalAuthorizedKey(sshPubKey)

	if err := os.WriteFile(pubPath, pubBytes, 0o644); err != nil {
		log.Debug().
			Str("publicKeyPath", pubPath).
			Err(err).
			Msg("failed to write public key")
		return "", fmt.Errorf("write public key: %w", err)
	}

	log.Debug().
		Str("publicKeyPath", pubPath).
		Msg("wrote public key to disk")

	return pubPath, nil
}
