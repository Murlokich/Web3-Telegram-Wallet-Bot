package repository

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/pkg/errors"
)

type EncryptedPostgres struct {
	encryptor encryption.Encryptor
	postgres  repository
	tracer    trace.Tracer
}

func New(tracer trace.Tracer, encryptor encryption.Encryptor, postgres repository) *EncryptedPostgres {
	return &EncryptedPostgres{encryptor: encryptor, postgres: postgres, tracer: tracer}
}

func (ep *EncryptedPostgres) AddNewAddress(ctx context.Context, userID int64) (*domain.AddressManagementData, error) {
	ctx, span := ep.tracer.Start(ctx, "AddNewAddress")
	defer span.End()
	resEncrypted, err := ep.postgres.AddNewAddress(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "failed to add new address to postgres")
		span.RecordError(err)
		return nil, err
	}
	res, err := resEncrypted.Decrypt(ctx, ep.encryptor)
	if err != nil {
		err = errors.Wrap(err, "failed to decrypt result")
		span.RecordError(err)
		return nil, err
	}
	return res, nil
}

func (ep *EncryptedPostgres) GetChangeLevelKey(
	ctx context.Context, userID int64) (*domain.AddressManagementData, error) {
	ctx, span := ep.tracer.Start(ctx, "ChangeLevelKey")
	defer span.End()
	resEncrypted, err := ep.postgres.GetChangeLevelKey(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "failed to get change level key from postgres")
		span.RecordError(err)
		return nil, err
	}
	res, err := resEncrypted.Decrypt(ctx, ep.encryptor)
	if err != nil {
		err = errors.Wrap(err, "failed to decrypt result")
		span.RecordError(err)
		return nil, err
	}
	return res, nil
}
func (ep *EncryptedPostgres) InsertWallet(ctx context.Context, wallet *domain.HDWallet) error {
	ctx, span := ep.tracer.Start(ctx, "InsertWallet")
	defer span.End()
	recordEncrypted, err := WalletEncryptedRecordFromDomain(ctx, wallet, ep.encryptor)
	if err != nil {
		err = errors.Wrap(err, "failed to encrypt record")
		span.RecordError(err)
		return err
	}
	if err = ep.postgres.InsertWallet(ctx, recordEncrypted); err != nil {
		err = errors.Wrap(err, "failed to insert wallet to postgres")
		span.RecordError(err)
		return err
	}
	return nil
}
