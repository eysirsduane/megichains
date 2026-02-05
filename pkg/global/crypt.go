package global

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"
)

var SecretKey = "_#D83^4*+678hfa^O#(39ru~!@#sf&aaaf*(&d*))U#p9{Dhus}&dhE*"

func stringToSecretKey(key string) [32]byte {
	hash := sha256.Sum256([]byte(key))
	return hash
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	return s
}

func EncryptEthPrivateKey(key string, serect string) (enrypted string, err error) {
	key = cleanString(key)
	if key == "" {
		err = errors.New("encrypted failed.")
		return
	}

	sbytes := stringToSecretKey(serect)
	pbytes, err := hex.DecodeString(key)
	if err != nil {
		return
	}
	if len(pbytes) == 0 {
		err = errors.New("encrypted failed..")
		return
	}

	var nonce [24]byte
	if _, err = rand.Read(nonce[:]); err != nil {
		return
	}

	ebytes := secretbox.Seal(nonce[:], pbytes, &nonce, &sbytes)

	return hex.EncodeToString(ebytes), nil
}

func DecryptEthPrivateKey(encrypted string, secret string) (decrypted string, err error) {
	sbytes := stringToSecretKey(secret)
	ebytes, err := hex.DecodeString(encrypted)
	if err != nil {
		return
	}

	if len(ebytes) < 24 {
		err = errors.New("decrypted failed.")
		return
	}

	var nonce [24]byte
	copy(nonce[:], ebytes[:24])
	cipherText := ebytes[24:]

	dbytes, ok := secretbox.Open(nil, cipherText, &nonce, &sbytes)
	if !ok {
		err = errors.New("decrypted failed.")
		return
	}

	decrypted = hex.EncodeToString(dbytes)

	return
}
