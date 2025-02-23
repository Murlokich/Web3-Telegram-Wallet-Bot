package handlers

import (
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v4"
)

type Handler func(ctx telebot.Context, deps *BotDependencies) error

type BotDependencies struct {
	Logger *logrus.Entry
	DB     *pgx.Conn
}

func (d *BotDependencies) UpdateLoggerFields(fields logrus.Fields) *BotDependencies {
	return &BotDependencies{
		Logger: d.Logger.WithFields(fields),
	}
}

func (d *BotDependencies) WrapHandler(handler Handler) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		logEntry := d.Logger.WithFields(logrus.Fields{
			"user_id": ctx.Sender().ID,
		})
		newDeps := &BotDependencies{
			Logger: logEntry,
			DB:     d.DB,
		}
		return handler(ctx, newDeps)
	}
}
