package adapter

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"context"
)

type HDWalletAdapter interface {
	GenerateHDWallet(ctx context.Context, userID int64) (*domain.HDWallet, string, error)
	DeriveWalletFromMnemonic(ctx context.Context, mnemonic string, userID int64) (*domain.HDWallet, error)
	GetAddress(ctx context.Context, changeLevelKeyBytes []byte, addressIndex uint32) (string, error)
}
