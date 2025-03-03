package handlers

import (
	"Web3-Telegram-Wallet-Bot/internal/db"
	"Web3-Telegram-Wallet-Bot/internal/wallet"
	"context"
	"fmt"

	"gopkg.in/telebot.v4"
)

const (
	internalErrorMessage = "Internal Server Error, please try again later"
	initialAddressIndex  = 0
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
	initialAddress, err := wallet.GetAddress(wlt.ChangeLevelKey, initialAddressIndex)
	if err != nil {
		dependencies.Logger.Error("Failed to get initial address: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	if err = db.InsertWallet(ctx, dependencies.DB, tgCtx.Sender().ID, mkEntry, clkEntry); err != nil {
		dependencies.Logger.Error("Failed to insert keys: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Please, remember your mnemonic:\n%s\n\nYour initial ETH address: %s",
		wlt.Mnemonic, initialAddress)

	return tgCtx.Send(message)
}

func AddNewAddress(tgCtx telebot.Context, dependencies *BotDependencies) error {
	ctx := context.Background()
	dependencies.Logger.Info("Creating Address")
	entry, lastAddressIndex, err := db.AddNewAddress(ctx, dependencies.DB, tgCtx.Sender().ID)
	if err != nil {
		dependencies.Logger.Error("Failed to get change level key: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	clk, err := dependencies.Encryptor.Decrypt(entry)
	if err != nil {
		dependencies.Logger.Error("Failed to decrypt entry: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	address, err := wallet.GetAddress(clk, lastAddressIndex)
	if err != nil {
		dependencies.Logger.Error("Failed to get address: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Your new ETH address is: %s", address)
	return tgCtx.Send(message)
}
