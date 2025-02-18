package handlers

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

type BotDependencies struct {
	Logger *logrus.Entry
}

func (d *BotDependencies) UpdateLoggerFields(fields logrus.Fields) *BotDependencies {
	return &BotDependencies{
		Logger: d.Logger.WithFields(fields),
	}
}

func (d *BotDependencies) WrapHandler(handler func(ctx telebot.Context, deps *BotDependencies) error) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		logEntry := d.Logger.WithFields(logrus.Fields{
			"user_id": ctx.Sender().ID,
		})
		newDeps := &BotDependencies{
			Logger: logEntry,
		}
		return handler(ctx, newDeps)
	}
}
