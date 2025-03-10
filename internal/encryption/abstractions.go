package encryption

type EncryptedEntry struct {
	Ciphertext []byte
	Nonce      []byte
}
