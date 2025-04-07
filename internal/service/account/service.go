package account

import (
	"Web3-Telegram-Wallet-Bot/internal/repository"
	"Web3-Telegram-Wallet-Bot/internal/service/adapter"

	"go.opentelemetry.io/otel/trace"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Logger   *logrus.Entry
	HDWallet adapter.HDWalletAdapter
	DB       repository.EncryptedRepository
	ETH      adapter.ETHAdapter
	tracer   trace.Tracer
}

func New(logger *logrus.Logger, hdWallet adapter.HDWalletAdapter, db repository.EncryptedRepository,
	eth adapter.ETHAdapter, tracer trace.Tracer) *Service {
	return &Service{
		Logger:   logger.WithField("service", "account"),
		HDWallet: hdWallet,
		DB:       db,
		ETH:      eth,
		tracer:   tracer,
	}
}
