package encryption

type Encryptor interface {
	Encrypt(plaintextBytes []byte) (*EncryptedEntry, error)
	Decrypt(entry *EncryptedEntry) ([]byte, error)
}
