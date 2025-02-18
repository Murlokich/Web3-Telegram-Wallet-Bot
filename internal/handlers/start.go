package handlers

import (
	"gopkg.in/telebot.v4"
)

const (
	internalErrorMessage = "Internal Server Error, please try again later"
)

func Start(ctx telebot.Context, dependencies *BotDependencies) error {
	dependencies.Logger.Info("Starting Telegram Bot")
	return ctx.Send("Hello, World!")
}
