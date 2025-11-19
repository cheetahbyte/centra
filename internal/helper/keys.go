package helper

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	xssh "golang.org/x/crypto/ssh"
)

func EnsureKeys(keysDir string) (string, error) {
	if err := os.MkdirAll(keysDir, 0o700); err != nil {
		return "", fmt.Errorf("create keys dir: %w", err)
	}

	privateKeyPath := filepath.Join(keysDir, "id_ed25519")
	publicKeyPath := filepath.Join(keysDir, "id_ed25519.pub")

	if _, err := os.Stat(privateKeyPath); err == nil {
		return publicKeyPath, nil
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("generate ed25519 key: %w", err)
	}

	privBlock, err := xssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return "", fmt.Errorf("marshal private key: %w", err)
	}

	privPEMBytes := pem.EncodeToMemory(privBlock)

	if err := os.WriteFile(privateKeyPath, privPEMBytes, 0o600); err != nil {
		return "", fmt.Errorf("write private key: %w", err)
	}

	sshPubKey, err := xssh.NewPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("create public key: %w", err)
	}

	pubBytes := xssh.MarshalAuthorizedKey(sshPubKey)

	if err := os.WriteFile(publicKeyPath, pubBytes, 0o644); err != nil {
		return "", fmt.Errorf("write public key: %w", err)
	}

	return publicKeyPath, nil
}
