package keystore

type SharedSecretGenerator func() ([]byte, error)

type KeyStore interface {
	EnsureSharedSecretSet(name string, generator SharedSecretGenerator) (SharedSecretSet, error)
}

type SharedSecretSet interface {
	EnsureSharedSecret() (SharedSecret, error)
}

type SharedSecret interface {
	SecretData() []byte
}
