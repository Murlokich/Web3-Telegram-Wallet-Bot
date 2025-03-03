package handlers

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

func RegisterBotHandlers(bot *telebot.Bot, dependencies *BotDependencies) {
	bot.Handle("/start", dependencies.UpdateLoggerFields(logrus.Fields{
		"command": "/start"},
	).WrapHandler(Start))

	bot.Handle("/register", dependencies.UpdateLoggerFields(logrus.Fields{
		"command": "/register"},
	).WrapHandler(CreateAccount))

	bot.Handle("/new_address", dependencies.UpdateLoggerFields(logrus.Fields{
		"command": "/new_address"},
	).WrapHandler(AddNewAddress))
}
