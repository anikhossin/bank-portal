package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// generateToken produces a random nonce signed with HMAC-SHA256, so the token
// is unguessable and tamper-evident. The stored/compared value is the full
// "nonce.signature" string — no secret needs to live in the database.
func generateToken(secret string) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	nonceHex := hex.EncodeToString(nonce)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(nonceHex))
	sig := hex.EncodeToString(mac.Sum(nil))
	return nonceHex + "." + sig, nil
}
