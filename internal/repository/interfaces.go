package repository

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"context"
)

type EncryptedRepository interface {
	AddNewAddress(ctx context.Context, userID int64) (*domain.AddressManagementData, error)
	GetChangeLevelKey(ctx context.Context, userID int64) (*domain.AddressManagementData, error)
	InsertWallet(ctx context.Context, record *domain.HDWallet) error
	UpdateCurrentAddress(ctx context.Context, userID int64, addressIndex uint32) error
}

type repository interface {
	AddNewAddress(ctx context.Context, userID int64) (*AddressManagementEncryptedData, error)
	GetChangeLevelKey(ctx context.Context, userID int64) (*AddressManagementEncryptedData, error)
	InsertWallet(ctx context.Context, record *WalletEncryptedRecord) error
	UpdateCurrentAddress(ctx context.Context, userID int64, addressIndex uint32) error
}
