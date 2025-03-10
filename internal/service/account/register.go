package account

import (
	"Web3-Telegram-Wallet-Bot/internal/service"
	"context"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
)

const (
	initialAddressIndex = 0
)

var (
	ErrInvalidMnemonic = errors.New("invalid mnemonic")
)

func (s *Service) CreateAccount(ctx context.Context, userID int64) (string, string, error) {
	s.Logger.Info("Creating Account")
	wlt, mnemonic, err := s.HDWallet.GenerateHDWallet(userID)
	if err != nil {
		s.Logger.Error("Failed to generate HDWallet: ", err)
		return "", "", service.ErrInternal
	}
	if err = s.DB.InsertWallet(ctx, wlt); err != nil {
		s.Logger.Error("Failed to insert wallet: ", err)
		return "", "", service.ErrInternal
	}
	initialAddress, err := s.HDWallet.GetAddress(wlt.AddressManagementData.ChangeLevelKey, initialAddressIndex)
	if err != nil {
		s.Logger.Error("Failed to get address: ", err)
		return "", "", service.ErrInternal
	}
	return mnemonic, initialAddress, nil
}

func (s *Service) MigrateAccount(ctx context.Context, mnemonic string, userID int64) (string, error) {
	s.Logger.Info("Migrating Account")
	if !bip39.IsMnemonicValid(mnemonic) {
		s.Logger.Error("Invalid mnemonic")
		return "", ErrInvalidMnemonic
	}
	wlt, err := s.HDWallet.DeriveWalletFromMnemonic(mnemonic, userID)
	if err != nil {
		s.Logger.Error("Failed to derive wallet from mnemonic: ", err)
		return "", service.ErrInternal
	}
	if err = s.DB.InsertWallet(ctx, wlt); err != nil {
		s.Logger.Error("Failed to insert wallet: ", err)
		return "", service.ErrInternal
	}
	initialAddress, err := s.HDWallet.GetAddress(wlt.AddressManagementData.ChangeLevelKey, initialAddressIndex)
	if err != nil {
		s.Logger.Error("Failed to get address: ", err)
		return "", service.ErrInternal
	}
	return initialAddress, nil
}
