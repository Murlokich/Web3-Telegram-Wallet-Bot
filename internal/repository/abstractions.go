package repository

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"Web3-Telegram-Wallet-Bot/internal/encryption"

	"github.com/pkg/errors"
)

func AddressManagementEncryptedDataFromDomain(
	data *domain.AddressManagementData, encryptor encryption.Encryptor) (*AddressManagementEncryptedData, error) {
	clkEntry, err := encryptor.Encrypt(data.ChangeLevelKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt change level key")
	}
	return &AddressManagementEncryptedData{ChangeLevelKey: *clkEntry, LastAddressIndex: data.LastAddressIndex}, nil
}

type WalletEncryptedRecord struct {
	UserID           int64
	MasterKey        *encryption.EncryptedEntry
	ChangeLevelKey   *encryption.EncryptedEntry
	LastAddressIndex uint32
}

func WalletEncryptedRecordFromDomain(
	wallet *domain.HDWallet, encryptor encryption.Encryptor) (*WalletEncryptedRecord, error) {
	mkEntry, err := encryptor.Encrypt(wallet.MasterKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt master key")
	}
	clkEntry, err := encryptor.Encrypt(wallet.AddressManagementData.ChangeLevelKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt change level key")
	}
	return &WalletEncryptedRecord{
		UserID:           wallet.UserID,
		MasterKey:        mkEntry,
		ChangeLevelKey:   clkEntry,
		LastAddressIndex: wallet.AddressManagementData.LastAddressIndex,
	}, nil
}

type AddressManagementEncryptedData struct {
	ChangeLevelKey   encryption.EncryptedEntry
	LastAddressIndex uint32
}

func (r *WalletEncryptedRecord) Decrypt(encryptor encryption.Encryptor) (*domain.HDWallet, error) {
	mk, err := encryptor.Decrypt(r.MasterKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt master key")
	}
	clk, err := encryptor.Decrypt(r.ChangeLevelKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt change level key")
	}
	return &domain.HDWallet{UserID: r.UserID, MasterKey: mk,
		AddressManagementData: &domain.AddressManagementData{ChangeLevelKey: clk, LastAddressIndex: r.LastAddressIndex}}, nil
}

func (r *AddressManagementEncryptedData) Decrypt(
	encryptor encryption.Encryptor) (*domain.AddressManagementData, error) {
	clk, err := encryptor.Decrypt(&r.ChangeLevelKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt change level key")
	}
	return &domain.AddressManagementData{ChangeLevelKey: clk, LastAddressIndex: r.LastAddressIndex}, nil
}
