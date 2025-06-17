package middlewares

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"echo-react-serve/config"
	"encoding/hex"
)

var (
	CookieSecret = []byte(config.Envs.App.CookieSecret) // Replace with a secure value
	ProtectAPI   = config.Envs.App.ProtectAPI
	RealBackend  = config.Envs.App.RealBackendUrl // Your real backend API
)

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func signValue(value string) string {
	h := hmac.New(sha256.New, CookieSecret)
	h.Write([]byte(value))
	return value + "|" + hex.EncodeToString(h.Sum(nil))
}

func validateSignedValue(signedValue string) (string, bool) {
	parts := []byte(signedValue)
	index := len(parts) - 65 // 64 hex + 1 separator
	if index <= 0 {
		return "", false
	}
	value := string(parts[:index])
	sig := string(parts[index+1:])

	expected := hmac.New(sha256.New, CookieSecret)
	expected.Write([]byte(value))
	expectedSig := hex.EncodeToString(expected.Sum(nil))

	return value, hmac.Equal([]byte(sig), []byte(expectedSig))
}
