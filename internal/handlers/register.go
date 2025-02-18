package handlers

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

func RegisterBotHandlers(bot *telebot.Bot, dependencies *BotDependencies) {
	bot.Handle("/start", dependencies.UpdateLoggerFields(logrus.Fields{
		"command": "/start"},
	).WrapHandler(Start))
}
