package crypt

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"passvault/pkg/log"
)

type CryptManager struct {
	logManager log.LogManagerInterface
	key        []byte
}

func NewCryptManager(logManager log.LogManagerInterface, key []byte) *CryptManager {
	cryptManager := &CryptManager{
		logManager: logManager,
		key:        key,
	}

	return cryptManager
}

func (cm *CryptManager) Encrypt(plaintext string) (*string, error) {
	c, err := aes.NewCipher(cm.key)
	if err != nil {
		//todo log
		return nil, err
	}

	plaintextBytes := []byte(plaintext)

	blockSize := c.BlockSize()
	padding := blockSize - (len(plaintextBytes) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintextBytes = append(plaintextBytes, padtext...)

	out := make([]byte, len(plaintextBytes))
	c.Encrypt(out, plaintextBytes)
	encToString := hex.EncodeToString(out)

	//todo log
	return &(encToString), nil
}

func (cm *CryptManager) Decrypt(encryptedHex string) (*string, error) {
	ciphertext, _ := hex.DecodeString(encryptedHex)
	c, err := aes.NewCipher(cm.key)
	if err != nil {
		//todo log
		return nil, err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	padLen := int(pt[len(pt)-1])
	s := string(pt[:len(pt)-padLen])

	//todo log
	return &s, nil
}
