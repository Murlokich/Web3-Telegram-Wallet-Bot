package aes

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"go.opentelemetry.io/otel/trace"
	"io"

	"github.com/pkg/errors"
)

const (
	encryptSpanName = "encrypt"
	decryptSpanName = "decrypt"
)

type Encryptor struct {
	gcm    cipher.AEAD
	tracer trace.Tracer
}

func New(tracer trace.Tracer, masterKeyB64 string) (*Encryptor, error) {
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
	return &Encryptor{gcm: gcm, tracer: tracer}, nil
}

func (e *Encryptor) Encrypt(ctx context.Context, plaintextBytes []byte) (*encryption.EncryptedEntry, error) {
	ctx, span := e.tracer.Start(ctx, encryptSpanName)
	defer span.End()
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		err = errors.Wrap(err, "failed to generate nonce")
		span.RecordError(err)
		return nil, err
	}
	ciphertext := e.gcm.Seal(nil, nonce, plaintextBytes, nil)

	return &encryption.EncryptedEntry{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	}, nil
}

func (e *Encryptor) Decrypt(ctx context.Context, entry *encryption.EncryptedEntry) ([]byte, error) {
	ctx, span := e.tracer.Start(ctx, decryptSpanName)
	defer span.End()
	decryptedBytes, err := e.gcm.Open(nil, entry.Nonce, entry.Ciphertext, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to decrypt ciphertext")
		span.RecordError(err)
		return nil, err
	}
	return decryptedBytes, nil
}
