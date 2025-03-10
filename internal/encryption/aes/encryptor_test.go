package aes_test

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"Web3-Telegram-Wallet-Bot/internal/encryption/aes"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	cycles = 3
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	encryptor, err := aes.New(base64.StdEncoding.EncodeToString(bytes))
	require.NoError(t, err)
	for step := range cycles {
		t.Run(strconv.Itoa(step), func(t *testing.T) {
			msg := make([]byte, 82)
			var entry *encryption.EncryptedEntry
			entry, err = encryptor.Encrypt(msg)
			require.NoError(t, err)
			require.NotNil(t, entry)
			require.NotNil(t, entry.Ciphertext)
			require.NotNil(t, entry.Nonce)
			require.NotEqual(t, msg, entry.Ciphertext)
			var decryptedMsg []byte
			decryptedMsg, err = encryptor.Decrypt(entry)
			require.NoError(t, err)
			require.NotNil(t, decryptedMsg)
			require.Equal(t, msg, decryptedMsg)
		})
	}
}

func TestEncryptDecrypt_SameText(t *testing.T) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	encryptor, err := aes.New(base64.StdEncoding.EncodeToString(bytes))
	require.NoError(t, err)
	msg := make([]byte, 82)
	entry1, err := encryptor.Encrypt(msg)
	require.NoError(t, err)
	entry2, err := encryptor.Encrypt(msg)
	require.NoError(t, err)
	require.NotEqual(t, entry1.Ciphertext, entry2.Ciphertext)
	require.NotEqual(t, entry1.Nonce, entry2.Nonce)
	msgOrig1, err := encryptor.Decrypt(entry1)
	require.NoError(t, err)
	msgOrig2, err := encryptor.Decrypt(entry2)
	require.NoError(t, err)
	require.Equal(t, msg, msgOrig1)
	require.Equal(t, msg, msgOrig2)
}

func TestEncryptDecrypt_Corruption(t *testing.T) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	encryptor, err := aes.New(base64.StdEncoding.EncodeToString(bytes))
	require.NoError(t, err)
	msg := make([]byte, 82)
	entry1, err := encryptor.Encrypt(msg)
	require.NoError(t, err)
	entry2, err := encryptor.Encrypt(msg)
	require.NoError(t, err)
	entry1.Ciphertext[0] ^= 0xff
	entry2.Nonce[0] ^= 0xff
	_, err = encryptor.Decrypt(entry1)
	require.Error(t, err)
	_, err = encryptor.Decrypt(entry2)
	require.Error(t, err)
}
