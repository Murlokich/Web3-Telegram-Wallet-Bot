package encryption

import "context"

type Encryptor interface {
	Encrypt(ctx context.Context, plaintextBytes []byte) (*EncryptedEntry, error)
	Decrypt(ctx context.Context, entry *EncryptedEntry) ([]byte, error)
}
