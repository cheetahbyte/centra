package helper

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheetahbyte/centra/internal/config"
	xssh "golang.org/x/crypto/ssh"
)

func EnsureKeys(keysDir string) (string, error) {
	if err := os.MkdirAll(keysDir, 0o700); err != nil {
		return "", fmt.Errorf("create keys dir: %w", err)
	}

	privateKeyPath := filepath.Join(keysDir, "id_ed25519")
	publicKeyPath := filepath.Join(keysDir, "id_ed25519.pub")

	configPrivKey := config.GetPrivateSSHKey()
	configPubKey := config.GetPublicSSHKey()

	if configPrivKey != "" {
		privKeyBytes := []byte(configPrivKey)
		if err := os.WriteFile(privateKeyPath, privKeyBytes, 0o600); err != nil {
			return "", fmt.Errorf("failed to write config private key: %w", err)
		}

		if configPubKey != "" {
			if err := os.WriteFile(publicKeyPath, []byte(configPubKey), 0o644); err != nil {
				return "", fmt.Errorf("failed to write config public key: %w", err)
			}
		} else {
			if err := deriveAndSavePublicKey(privKeyBytes, publicKeyPath); err != nil {
				return "", fmt.Errorf("failed to derive public key from config: %w", err)
			}
		}

		return publicKeyPath, nil
	}

	if _, err := os.Stat(privateKeyPath); err == nil {
		if _, err := os.Stat(publicKeyPath); err == nil {
			return publicKeyPath, nil
		}

		existingPrivBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return "", fmt.Errorf("read existing private key: %w", err)
		}
		if err := deriveAndSavePublicKey(existingPrivBytes, publicKeyPath); err != nil {
			return "", fmt.Errorf("derive public key from existing file: %w", err)
		}
		return publicKeyPath, nil
	}

	return generateAndSaveKeys(privateKeyPath, publicKeyPath)
}

func deriveAndSavePublicKey(privKeyBytes []byte, pubPath string) error {
	signer, err := xssh.ParsePrivateKey(privKeyBytes)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	sshPubKey := signer.PublicKey()

	pubBytes := xssh.MarshalAuthorizedKey(sshPubKey)

	if err := os.WriteFile(pubPath, pubBytes, 0o644); err != nil {
		return fmt.Errorf("write derived public key: %w", err)
	}

	return nil
}

func generateAndSaveKeys(privPath, pubPath string) (string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("generate ed25519 key: %w", err)
	}

	privBlock, err := xssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return "", fmt.Errorf("marshal private key: %w", err)
	}

	privPEMBytes := pem.EncodeToMemory(privBlock)

	if err := os.WriteFile(privPath, privPEMBytes, 0o600); err != nil {
		return "", fmt.Errorf("write private key: %w", err)
	}

	sshPubKey, err := xssh.NewPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("create public key: %w", err)
	}

	pubBytes := xssh.MarshalAuthorizedKey(sshPubKey)

	if err := os.WriteFile(pubPath, pubBytes, 0o644); err != nil {
		return "", fmt.Errorf("write public key: %w", err)
	}

	return pubPath, nil
}
