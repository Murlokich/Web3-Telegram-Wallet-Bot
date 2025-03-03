package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
)

type Encryptor struct {
	gcm cipher.AEAD
}

type EncryptedEntry struct {
	Ciphertext []byte
	Nonce      []byte
}

func NewEncryptor(masterKeyB64 string) (*Encryptor, error) {
	masterKey, err := base64.StdEncoding.DecodeString(masterKeyB64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode master key")
	}
	cipherBlock, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher block")
	}
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create GCM")
	}
	return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(plaintextBytes []byte) (*EncryptedEntry, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "failed to generate nonce")
	}
	ciphertext := e.gcm.Seal(nil, nonce, plaintextBytes, nil)

	return &EncryptedEntry{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	}, nil
}

func (e *Encryptor) Decrypt(entry *EncryptedEntry) ([]byte, error) {
	decryptedBytes, err := e.gcm.Open(nil, entry.Nonce, entry.Ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt ciphertext")
	}
	return decryptedBytes, nil
}
