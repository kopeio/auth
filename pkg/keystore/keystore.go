package keystore

//type sharedSecretGenerator func() ([]byte, error)

type KeyStore interface {
	KeySet(keyname string) (KeySet, error)
}

type KeySet interface {
	Decrypt(ciphertext []byte) ([]byte, error)
	Encrypt(plaintext []byte) ([]byte, error)
}

//type SharedSecretSet interface {
//	EnsureSharedSecret() (SharedSecret, error)
//}

//type SharedSecret interface {
//	SecretData() []byte
//}

//type Key interface {
//	Id() int32
//	//Secret() []byte
//}
