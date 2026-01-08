package helper

import (
	"strings"

	"golang.org/x/crypto/ssh"
)

func PublicKeyToString(pub ssh.PublicKey) string {
	// MarshalAuthorizedKey returns bytes in the format: "type base64-key\n"
	pubBytes := ssh.MarshalAuthorizedKey(pub)

	// Convert to string and trim the trailing newline
	return strings.TrimSpace(string(pubBytes))
}
