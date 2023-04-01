package crypt

type CryptManagerInterface interface {
	Encrypt(string) (*string, error)
	Decrypt(string) (*string, error)
}
