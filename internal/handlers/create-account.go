package handlers

import (
	"Web3-Telegram-Wallet-Bot/internal/db"
	"Web3-Telegram-Wallet-Bot/internal/wallet"
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"

	"gopkg.in/telebot.v4"
)

const (
	invalidMnemonic      = "Mnemonic you have provided is invalid. Please check your input and try again."
	internalErrorMessage = "Internal Server Error, please try again later"
	invalidIndex         = "Index you have provided is invalid. Please check your input and try again."
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
	initialAddress, err := saveWallet(ctx, dependencies, wlt, tgCtx.Sender().ID)
	if err != nil {
		dependencies.Logger.Error("Failed to save wallet: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Please, remember your mnemonic:\n%s\n\nYour initial ETH address: %s",
		wlt.Mnemonic, initialAddress)
	return tgCtx.Send(message)
}

func MigrateAccount(tgCtx telebot.Context, dependencies *BotDependencies) error {
	ctx := context.Background()
	dependencies.Logger.Info("Migrating Account")
	mnemonic := tgCtx.Data()
	if !bip39.IsMnemonicValid(mnemonic) {
		dependencies.Logger.Error("Invalid mnemonic")
		return tgCtx.Send(invalidMnemonic)
	}
	wlt, err := wallet.DeriveWalletFromMnemonic(mnemonic)
	if err != nil {
		dependencies.Logger.Error("Failed to derive wallet from mnemonic: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	initialAddress, err := saveWallet(ctx, dependencies, wlt, tgCtx.Sender().ID)
	if err != nil {
		dependencies.Logger.Error("Failed to save wallet: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Your initial ETH address: %s", initialAddress)
	return tgCtx.Send(message)
}

func saveWallet(ctx context.Context, dependencies *BotDependencies, wlt *wallet.HDWallet, userID int64) (string, error) {
	mkEntry, err := dependencies.Encryptor.Encrypt(wlt.MasterKey)
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt master key")
	}
	clkEntry, err := dependencies.Encryptor.Encrypt(wlt.ChangeLevelKey)
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt change level key")
	}
	initialAddress, err := wallet.GetAddress(wlt.ChangeLevelKey, initialAddressIndex)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get initial address")
	}
	if err = db.InsertWallet(ctx, dependencies.DB, userID, mkEntry, clkEntry); err != nil {
		return "", errors.Wrap(err, "Failed to insert wallet")
	}
	return initialAddress, nil
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

func GetBalance(tgCtx telebot.Context, dependencies *BotDependencies) error {
	ctx := context.Background()
	dependencies.Logger.Info("Getting Balance")
	entry, lindex, err := db.GetChangeLevelKey(ctx, dependencies.DB, tgCtx.Sender().ID)
	if err != nil {
		dependencies.Logger.Error("Failed to get change level key: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	index, err := strconv.ParseInt(tgCtx.Data(), 10, 4)
	if err != nil {
		dependencies.Logger.Error("Failed to parse index: ", err)
		return tgCtx.Send(invalidIndex)
	}
	indexu32 := uint32(index)
	if index < 0 || indexu32 > lindex {
		dependencies.Logger.Error("Invalid index")
		return tgCtx.Send(invalidIndex)
	}
	clk, err := dependencies.Encryptor.Decrypt(entry)
	if err != nil {
		dependencies.Logger.Error("Failed to decrypt entry: ", err)
		return tgCtx.Send(internalErrorMessage)
	}
	addr, err := wallet.GetAddress(clk, indexu32)
	if err != nil {
		dependencies.Logger.Error("Failed to get address: ", err)
		return tgCtx.Send(internalErrorMessage)
	}

}
