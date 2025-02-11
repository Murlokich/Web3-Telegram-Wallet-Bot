package commands

import "gopkg.in/telebot.v4"

func Start(ctx telebot.Context) error {
	return ctx.Send("Hello, World!")
}
