package handlers

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"
)

type DBProvider interface {
	AddNewAddress(ctx context.Context, userID int64) (*encryption.EncryptedEntry, uint32, error)
	GetChangeLevelKey(ctx context.Context, userID int64) (*encryption.EncryptedEntry, uint32, error)
	InsertWallet(ctx context.Context, userID int64,
		mkEntry *encryption.EncryptedEntry, clkEntry *encryption.EncryptedEntry) error
}
