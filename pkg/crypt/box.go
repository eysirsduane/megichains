package crypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

var PrivateKeySecretPassword = "Megichains-*fljk234fl*Fj10;fi023897#$*49*%^37O(5BJF9833h_*"
var PrivateKeySecretSalt = "_#D83^4*+678hfa^O#(39ru~!@#sf&aaaf*(&d*))U#p9{Dhus}&dhE*"

const (
	keyLength = 32
	nonceSize = 24
)

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

func DeriveKey(password string, salt []byte) ([32]byte, error) {
	var key [32]byte

	dk, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, keyLength)
	if err != nil {
		return key, err
	}

	copy(key[:], dk)
	return key, nil
}

func Encrypt(plain, password, salt string) (string, error) {
	key, err := DeriveKey(password, []byte(salt))
	if err != nil {
		return "", err
	}

	var nonce [nonceSize]byte

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	encrypted := secretbox.Seal(nonce[:], []byte(plain), &nonce, &key)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(cipherText, password, salt string) (string, error) {
	key, err := DeriveKey(password, []byte(salt))
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	var nonce [nonceSize]byte
	copy(nonce[:], data[:nonceSize])

	decrypted, ok := secretbox.Open(nil, data[nonceSize:], &nonce, &key)
	if !ok {
		return "", fmt.Errorf("decryption failed")
	}

	return string(decrypted), nil
}
