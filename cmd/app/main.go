package main

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process config: %v", err)
		panic(err)
	}
}
