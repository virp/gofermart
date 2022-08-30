package web

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

func SetEncryptedCookie(w http.ResponseWriter, name string, value string, secret string) error {

	aesGSM, err := getAES(secret)
	if err != nil {
		return fmt.Errorf("get AES: %w", err)
	}

	nonce, err := generateRandom(aesGSM.NonceSize())
	if err != nil {
		return fmt.Errorf("generate random nonce: %w", err)
	}

	v := aesGSM.Seal(nil, nonce, []byte(value), nil)
	v = append(v, nonce...)

	c := &http.Cookie{
		Name:  name,
		Value: hex.EncodeToString(v),
	}

	http.SetCookie(w, c)

	return nil
}

func GetEncryptedCookie(r *http.Request, name string, secret string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("get cookie from request: %w", err)
	}

	aesGSM, err := getAES(secret)
	if err != nil {
		return "", fmt.Errorf("get AES: %w", err)
	}

	v, err := hex.DecodeString(c.Value)
	if err != nil {
		return "", fmt.Errorf("decode cookie value: %w", err)
	}

	nonce := v[len(v)-aesGSM.NonceSize():]
	value, err := aesGSM.Open(nil, nonce, v[:len(v)-aesGSM.NonceSize()], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt cookie value: %w", err)
	}

	return string(value), nil
}

func getAES(secret string) (cipher.AEAD, error) {
	key := sha256.Sum256([]byte(secret))
	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("aesblock: %w", err)
	}
	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("aesgcm: %w", err)
	}

	return aesGCM, nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
