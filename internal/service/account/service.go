package account

import (
	"Web3-Telegram-Wallet-Bot/internal/repository"
	"Web3-Telegram-Wallet-Bot/internal/service/adapter"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Logger   *logrus.Entry
	HDWallet adapter.HDWalletAdapter
	DB       repository.EncryptedRepository
}

func New(logger *logrus.Logger, hdWallet adapter.HDWalletAdapter, db repository.EncryptedRepository) *Service {
	return &Service{
		Logger:   logger.WithField("service", "account"),
		HDWallet: hdWallet,
		DB:       db,
	}
}
