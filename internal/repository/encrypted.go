package repository

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"

	"github.com/pkg/errors"
)

type EncryptedPostgres struct {
	encryptor encryption.Encryptor
	postgres  repository
}

func New(encryptor encryption.Encryptor, postgres repository) *EncryptedPostgres {
	return &EncryptedPostgres{encryptor: encryptor, postgres: postgres}
}

func (ep *EncryptedPostgres) AddNewAddress(ctx context.Context, userID int64) (*domain.AddressManagementData, error) {
	resEncrypted, err := ep.postgres.AddNewAddress(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new address to postgres")
	}
	res, err := resEncrypted.Decrypt(ep.encryptor)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt result")
	}
	return res, nil
}

func (ep *EncryptedPostgres) GetChangeLevelKey(
	ctx context.Context, userID int64) (*domain.AddressManagementData, error) {
	resEncrypted, err := ep.postgres.GetChangeLevelKey(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get change level key from postgres")
	}
	res, err := resEncrypted.Decrypt(ep.encryptor)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt result")
	}
	return res, nil
}
func (ep *EncryptedPostgres) InsertWallet(ctx context.Context, wallet *domain.HDWallet) error {
	recordEncrypted, err := WalletEncryptedRecordFromDomain(wallet, ep.encryptor)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt record")
	}
	if err = ep.postgres.InsertWallet(ctx, recordEncrypted); err != nil {
		return errors.Wrap(err, "failed to insert wallet to postgres")
	}
	return nil
}
