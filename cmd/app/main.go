package main

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"Web3-Telegram-Wallet-Bot/internal/controller/telegram"
	"Web3-Telegram-Wallet-Bot/internal/encryption/aes"
	"Web3-Telegram-Wallet-Bot/internal/repository"
	postgres2 "Web3-Telegram-Wallet-Bot/internal/repository/postgres"
	"Web3-Telegram-Wallet-Bot/internal/service/account"
	"Web3-Telegram-Wallet-Bot/internal/service/adapter/eth/infura"
	"Web3-Telegram-Wallet-Bot/internal/service/adapter/wallet/bip32adapter"
	"Web3-Telegram-Wallet-Bot/internal/tracing"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

func runMigrations(dbConfig *config.DBConfig) error {
	m, err := migrate.New("file://migrations", dbConfig.URL)
	if err != nil {
		return errors.Wrap(err, "failed to create migration")
	}
	err = m.Migrate(dbConfig.MigrationVersion)
	if err != nil {
		return errors.Wrap(err, "failed to migrate database")
	}
	errSrc, errDB := m.Close()
	if errSrc != nil {
		return errors.Wrap(err, "failed to close migration source")
	}
	if errDB != nil {
		return errors.Wrap(err, "failed to close database connection")
	}
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Errorf("failed to process config: %v", err)
		return
	}

	tracerProvider, err := tracing.NewTracerProvider(ctx, &cfg.Tracing)
	if err != nil {
		log.Errorf("failed to initialize tracer provider: %v", err)
		return
	}

	botSettings := telebot.Settings{
		Token:  cfg.TelegramBotConfig.Token,
		Poller: &telebot.LongPoller{Timeout: time.Duration(cfg.TelegramBotConfig.Timeout) * time.Second},
	}

	bot, err := telebot.NewBot(botSettings)
	if err != nil {
		log.Errorf("failed to create bot: %v", err)
		return
	}

	if err = runMigrations(&cfg.DBConfig); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Errorf("failed to run migrations: %v", err)
		return
	}

	encryptor, err := aes.New(tracerProvider.Tracer("encryptor"), cfg.Encryption.MasterKey)
	if err != nil {
		log.Errorf("failed to create encryptor: %v", err)
		return
	}

	postgresClient, err := postgres2.New(ctx, tracerProvider.Tracer("postgres"), &cfg.DBConfig)
	if err != nil {
		log.Errorf("failed to create postgres client: %v", err)
		return
	}
	encryptedPostgres := repository.New(tracerProvider.Tracer("encryptedPostgres"), encryptor, postgresClient)
	hdWalletAdapter := bip32adapter.New(tracerProvider.Tracer("bip32adapter"))
	ethProvider := infura.New(&cfg.Infura)

	accountService := account.New(
		log, hdWalletAdapter, encryptedPostgres, ethProvider, tracerProvider.Tracer("account-service"),
	)

	services := &telegram.BotServices{
		Logger:         log,
		AccountService: accountService,
	}
	telegram.RegisterBotHandlers(bot, services)

	go bot.Start()

	// graceful shutdown
	<-ctx.Done()

	stop()
	log.Infoln("shutting down gracefully")
	if err = postgresClient.Close(); err != nil {
		log.Errorf("failed to close connection to database: %v", err)
	} else {
		log.Info("connection to database closed")
	}
	log.Info("shutdown completed")
}
