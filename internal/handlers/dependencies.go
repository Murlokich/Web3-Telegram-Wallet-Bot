package handlers

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

type Handler func(ctx telebot.Context, deps *BotDependencies) error

type BotDependencies struct {
	Logger    *logrus.Entry
	DB        DBProvider
	Encryptor *encryption.Encryptor
}

func (d *BotDependencies) UpdateLoggerFields(fields logrus.Fields) *BotDependencies {
	newDependencies := *d
	newDependencies.Logger = d.Logger.WithFields(fields)
	return &newDependencies
}

func (d *BotDependencies) WrapHandler(handler Handler) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		return handler(ctx, d.UpdateLoggerFields(logrus.Fields{"user_id": ctx.Sender().ID}))
	}
}
