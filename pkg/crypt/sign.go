package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSHA256(message, secret string) (enc string) {
	key := []byte(secret)
	msg := []byte(message)

	h := hmac.New(sha256.New, key)

	h.Write(msg)

	enc = hex.EncodeToString(h.Sum(nil))

	return
}

func VerifySignature(csign string, message []byte, secret string) (ok bool) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)

	ssign := hex.EncodeToString(mac.Sum(nil))

	if hmac.Equal([]byte(ssign), []byte(csign)) {
		ok = true
		return
	}

	return
}
