package telegram

import (
	"Web3-Telegram-Wallet-Bot/internal/service"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

type Handler func(ctx telebot.Context, deps *BotServices) error

type BotServices struct {
	Logger         *logrus.Logger
	AccountService service.AccountService
}

func (d *BotServices) WrapHandler(handler Handler) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		return handler(ctx, d)
	}
}
