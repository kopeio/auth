package keystore

import (
	crypto_rand "crypto/rand"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/nacl/secretbox"

	"kope.io/auth/pkg/keystore/pb"
)

type secretboxKey struct {
	data pb.KeyData
}

func (k *secretboxKey) Id() int32 {
	return k.data.Id
}

func (k *secretboxKey) encrypt(plaintext []byte) ([]byte, error) {
	// From the example in the secretbox docs:
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		return nil, fmt.Errorf("error reading random data: %w", err)
	}

	secretKey := k.data.Secret
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("expected 32 byte key, was %d", len(secretKey))
	}

	var secretKeyArray [32]byte
	copy(secretKeyArray[:], secretKey[:32])

	ciphertext := secretbox.Seal(nil, plaintext, &nonce, &secretKeyArray)

	encrypted := &pb.EncryptedData{
		EncryptionMethod: pb.EncryptionMethod_ENCRYPTIONMETHOD_SECRETBOX,
		KeyId:            k.Id(),
		Nonce:            nonce[:],
		Ciphertext:       ciphertext,
	}

	encryptedBytes, err := proto.Marshal(encrypted)
	if err != nil {
		return nil, fmt.Errorf("error serializing data: %v", err)
	}

	return encryptedBytes, nil
}

func (k *secretboxKey) decrypt(encryptedData *pb.EncryptedData) ([]byte, error) {
	if len(encryptedData.Nonce) != 24 {
		return nil, fmt.Errorf("invalid nonce data")
	}

	secretKey := k.data.Secret
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("expected 32 byte key, was %d", len(secretKey))
	}

	var nonceArray [24]byte
	copy(nonceArray[:], encryptedData.Nonce)

	var secretKeyArray [32]byte
	copy(secretKeyArray[:], secretKey[:32])

	plaintext, ok := secretbox.Open(nil, encryptedData.Ciphertext, &nonceArray, &secretKeyArray)
	if !ok {
		return nil, fmt.Errorf("encrypted data not valid")
	}

	return plaintext, nil
}
