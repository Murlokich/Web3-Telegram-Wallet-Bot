package adapter

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
)

type HDWalletAdapter interface {
	GenerateHDWallet(userID int64) (*domain.HDWallet, string, error)
	DeriveWalletFromMnemonic(mnemonic string, userID int64) (*domain.HDWallet, error)
	GetAddress(changeLevelKeyBytes []byte, addressIndex uint32) (string, error)
}
