package telegram

import (
	"gopkg.in/telebot.v4"
)

func RegisterBotHandlers(bot *telebot.Bot, dependencies *BotServices) {
	bot.Handle("/start", Start)
	bot.Handle("/register", dependencies.WrapHandler(CreateAccount))
	bot.Handle("/new_address", dependencies.WrapHandler(AddNewAddress))
	bot.Handle("/migrate", dependencies.WrapHandler(MigrateAccount))
}
