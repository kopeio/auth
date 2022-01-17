package keystore

import "context"

type KeyStore interface {
	KeySet(ctx context.Context, keyname string) (KeySet, error)
}

type KeySet interface {
	Decrypt(ciphertext []byte) ([]byte, error)
	Encrypt(plaintext []byte) ([]byte, error)
}
