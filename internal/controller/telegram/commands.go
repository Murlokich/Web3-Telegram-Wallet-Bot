package telegram

import (
	"Web3-Telegram-Wallet-Bot/internal/service/account"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/telebot.v4"
)

const (
	invalidMnemonic      = "Mnemonic you have provided is invalid. Please check your input and try again."
	internalErrorMessage = "Internal Server Error, please try again later"
)

func Start(ctx telebot.Context) error {
	return ctx.Send("Hello, World!")
}

func CreateAccount(tgCtx telebot.Context, dependencies *BotServices) error {
	mnemonic, address, err := dependencies.AccountService.CreateAccount(context.Background(), tgCtx.Sender().ID)
	if err != nil {
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Please, remember your mnemonic:\n%s\n\nYour initial ETH address: %s",
		mnemonic, address)
	return tgCtx.Send(message)
}

func MigrateAccount(tgCtx telebot.Context, dependencies *BotServices) error {
	address, err := dependencies.AccountService.MigrateAccount(context.Background(), tgCtx.Data(), tgCtx.Sender().ID)
	if err != nil {
		if errors.Is(err, account.ErrInvalidMnemonic) {
			return tgCtx.Send(invalidMnemonic)
		}
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Your initial ETH address: %s", address)
	return tgCtx.Send(message)
}

func AddNewAddress(tgCtx telebot.Context, dependencies *BotServices) error {
	address, err := dependencies.AccountService.AddNewAddress(context.Background(), tgCtx.Sender().ID)
	if err != nil {
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Your new ETH address is: %s", address)
	return tgCtx.Send(message)
}

func GetBalance(tgCtx telebot.Context, dependencies *BotServices) error {
	ethBalance, err := dependencies.AccountService.GetBalance(context.Background(), tgCtx.Data())
	if err != nil {
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Balance for provided address is: %s ETH", ethBalance)
	return tgCtx.Send(message)
}
