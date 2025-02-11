package main

import (
	"Web3-Telegram-Wallet-Bot/internal/commands"
	"Web3-Telegram-Wallet-Bot/internal/config"
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logrus.New()
	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process config: %v", err)
		stop()
	}

	botSettings := telebot.Settings{
		Token:  cfg.TelegramBotConfig.Token,
		Poller: &telebot.LongPoller{Timeout: time.Duration(cfg.TelegramBotConfig.Timeout) * time.Second},
	}
	bot, err := telebot.NewBot(botSettings)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
		stop()
	}

	bot.Handle("/start", commands.Start)

	go bot.Start()

	// graceful shutdown
	<-ctx.Done()
	log.Infoln("shutting down gracefully")
}
