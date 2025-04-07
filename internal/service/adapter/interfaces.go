package adapter

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"context"
	"math/big"
)

type HDWalletAdapter interface {
	GenerateHDWallet(ctx context.Context, userID int64) (*domain.HDWallet, string, error)
	DeriveWalletFromMnemonic(ctx context.Context, mnemonic string, userID int64) (*domain.HDWallet, error)
	GetAddress(ctx context.Context, changeLevelKeyBytes []byte, addressIndex uint32) (string, error)
}

type ETHAdapter interface {
	// GetBalance returns balance in wei
	GetBalance(ctx context.Context, address string) (*big.Int, error)
}
