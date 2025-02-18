package main

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"Web3-Telegram-Wallet-Bot/internal/handlers"
	"context"
	"os/signal"
	"syscall"
	"time"

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

	dependencies := &handlers.BotDependencies{Logger: log.WithFields(logrus.Fields{})}
	handlers.RegisterBotHandlers(bot, dependencies)

	go bot.Start()

	// graceful shutdown
	<-ctx.Done()
	log.Infoln("shutting down gracefully")
}
