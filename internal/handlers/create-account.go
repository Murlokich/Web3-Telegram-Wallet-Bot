package handlers

import (
	"Web3-Telegram-Wallet-Bot/internal/db"
	"Web3-Telegram-Wallet-Bot/internal/wallet"
	"context"

	"gopkg.in/telebot.v4"
)

const (
	internalErrorMessage = "Internal Server Error, please try again later"
)

func Start(ctx telebot.Context, dependencies *BotDependencies) error {
	dependencies.Logger.Info("Starting Telegram Bot")
	return ctx.Send("Hello, World!")
}

func CreateAccount(tgCtx telebot.Context, dependencies *BotDependencies) error {
	ctx := context.Background()
	dependencies.Logger.Info("Creating Account")
	wlt, err := wallet.GenerateHDWallet()
	if err != nil {
		dependencies.Logger.Error("Failed to generate HDWallet: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	mkEntry, err := dependencies.Encryptor.Encrypt(wlt.MasterKey)
	if err != nil {
		dependencies.Logger.Error("Failed to encrypt master key: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	clkEntry, err := dependencies.Encryptor.Encrypt(wlt.ChangeLevelKey)
	if err != nil {
		dependencies.Logger.Error("Failed to encrypt change level key: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	if err = db.InsertKeys(ctx, dependencies.DB, tgCtx.Sender().ID, mkEntry, clkEntry); err != nil {
		dependencies.Logger.Error("Failed to insert keys: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	return tgCtx.Send(wlt.Mnemonic)
}
