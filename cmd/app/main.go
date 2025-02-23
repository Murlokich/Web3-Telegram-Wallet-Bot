package main

import (
	"Web3-Telegram-Wallet-Bot/db"
	"Web3-Telegram-Wallet-Bot/internal/config"
	"Web3-Telegram-Wallet-Bot/internal/handlers"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

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

	botSettings := telebot.Settings{
		Token:  cfg.TelegramBotConfig.Token,
		Poller: &telebot.LongPoller{Timeout: time.Duration(cfg.TelegramBotConfig.Timeout) * time.Second},
	}

	bot, err := telebot.NewBot(botSettings)
	if err != nil {
		log.Errorf("failed to create bot: %v", err)
		return
	}

	if err = db.RunMigrations(&cfg.DBConfig); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Errorf("failed to run migrations: %v", err)
		return
	}

	conn, err := pgx.Connect(ctx, cfg.DBConfig.URL)
	if err != nil {
		log.Errorf("failed to connect to database: %v", err)
		return
	}

	dependencies := &handlers.BotDependencies{
		Logger: log.WithFields(logrus.Fields{}),
		DB:     conn,
	}
	handlers.RegisterBotHandlers(bot, dependencies)

	go bot.Start()

	// graceful shutdown
	<-ctx.Done()

	stop()
	log.Infoln("shutting down gracefully")
	if err = conn.Close(context.Background()); err != nil {
		log.Errorf("failed to close connection to database: %v", err)
	} else {
		log.Info("connection to database closed")
	}
	log.Info("shutdown completed")
}
