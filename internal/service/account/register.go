package account

import (
	"Web3-Telegram-Wallet-Bot/internal/service"
	"context"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
)

const (
	initialAddressIndex = 0

	createAccountSpanName  = "CreateAccount"
	migrateAccountSpanName = "MigrateAccount"
)

var (
	ErrInvalidMnemonic = errors.New("invalid mnemonic")
)

func (s *Service) CreateAccount(ctx context.Context, userID int64) (string, string, error) {
	ctx, span := s.tracer.Start(ctx, createAccountSpanName)
	defer span.End()
	s.Logger.Info("Creating Account")
	wlt, mnemonic, err := s.HDWallet.GenerateHDWallet(ctx, userID)
	if err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to generate HDWallet: ", err)
		return "", "", service.ErrInternal
	}
	if err = s.DB.InsertWallet(ctx, wlt); err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to insert wallet: ", err)
		return "", "", service.ErrInternal
	}
	initialAddress, err := s.HDWallet.GetAddress(ctx, wlt.AddressManagementData.ChangeLevelKey, initialAddressIndex)
	if err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to get address: ", err)
		return "", "", service.ErrInternal
	}
	return mnemonic, initialAddress, nil
}

func (s *Service) MigrateAccount(ctx context.Context, mnemonic string, userID int64) (string, error) {
	ctx, span := s.tracer.Start(ctx, migrateAccountSpanName)
	defer span.End()
	s.Logger.Info("Migrating Account")
	if !bip39.IsMnemonicValid(mnemonic) {
		span.RecordError(ErrInvalidMnemonic)
		s.Logger.Error("Invalid mnemonic")
		return "", ErrInvalidMnemonic
	}
	wlt, err := s.HDWallet.DeriveWalletFromMnemonic(ctx, mnemonic, userID)
	if err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to derive wallet from mnemonic: ", err)
		return "", service.ErrInternal
	}
	if err = s.DB.InsertWallet(ctx, wlt); err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to insert wallet: ", err)
		return "", service.ErrInternal
	}
	initialAddress, err := s.HDWallet.GetAddress(ctx, wlt.AddressManagementData.ChangeLevelKey, initialAddressIndex)
	if err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to get address: ", err)
		return "", service.ErrInternal
	}
	return initialAddress, nil
}
