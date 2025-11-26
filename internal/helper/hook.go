package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

func VerifySignature(body []byte, signatureHeader, secret string) error {
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return errors.New("signature header is missing 'sha256=' prefix")
	}

	hexSig := strings.TrimPrefix(signatureHeader, "sha256=")
	receivedSig, err := hex.DecodeString(hexSig)
	if err != nil {
		return fmt.Errorf("failed to decode signature hex: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(expectedSig, receivedSig) {
		return errors.New("signature does not match expected hmac")
	}

	return nil
}
