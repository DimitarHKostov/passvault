package crypt

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
)

var (
	key          = []byte("this is secret key enough 32 bit")
	cryptManager *CryptManager
)

type CryptManager struct{}

func Get() *CryptManager {
	if cryptManager == nil {
		cryptManager = &CryptManager{}
	}

	return cryptManager
}

func (cm *CryptManager) Encrypt(plaintext string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)

	blockSize := c.BlockSize()
	padding := blockSize - (len(plaintextBytes) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintextBytes = append(plaintextBytes, padtext...)

	out := make([]byte, len(plaintextBytes))
	c.Encrypt(out, plaintextBytes)
	return hex.EncodeToString(out), nil
}

func (cm *CryptManager) Decrypt(encryptedHex string) (string, error) {
	ciphertext, _ := hex.DecodeString(encryptedHex)
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	padLen := int(pt[len(pt)-1])
	s := string(pt[:len(pt)-padLen])

	return s, nil
}
