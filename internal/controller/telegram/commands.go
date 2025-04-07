package telegram

import (
	"Web3-Telegram-Wallet-Bot/internal/service/account"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/telebot.v4"
)

const (
	invalidMnemonic      = "Mnemonic you have provided is invalid. Please check your input and try again."
	internalErrorMessage = "Internal Server Error, please try again later"

	addressOptionButtonEndpoint = "chooseAddress"

	addressButtonPayloadSegments = 2
)

func Start(ctx telebot.Context) error {
	return ctx.Send("Hello, World!", telebot.RemoveKeyboard)
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

func SwitchCurrentAddress(tgCtx telebot.Context, dependencies *BotServices) error {
	userID := tgCtx.Sender().ID
	addresses, err := dependencies.AccountService.GetAllAddresses(context.Background(), userID)
	if err != nil {
		return tgCtx.Send(internalErrorMessage)
	}
	btns := make([][]telebot.InlineButton, len(addresses))
	for i, address := range addresses {
		btns[i] = []telebot.InlineButton{{
			Unique: addressOptionButtonEndpoint, // All buttons share this one handler
			Data:   fmt.Sprintf("%d:%s", i, address),
			Text:   address,
		}}
	}
	return tgCtx.Send("Switch to address:", &telebot.ReplyMarkup{
		InlineKeyboard:  btns,
		OneTimeKeyboard: true,
		ForceReply:      true,
		ResizeKeyboard:  true,
	}, telebot.Protected, // Protected forbids forwarding/saving this message, so user cannot give us another user ID
	)
}

func ApplyAddressSwitch(tgCtx telebot.Context, dependencies *BotServices) error {
	data := tgCtx.Data()

	splitData := strings.Split(data, ":")
	if len(splitData) != addressButtonPayloadSegments {
		dependencies.Logger.Errorf(
			"ApplyAddressSwitch: address button data is split into %d strings, expected %d",
			len(splitData), addressButtonPayloadSegments)
		return tgCtx.Send(internalErrorMessage)
	}
	addressIndex, err := strconv.Atoi(splitData[0])
	if err != nil {
		dependencies.Logger.Error("ApplyAddressSwitch: address index is not a number")
		return tgCtx.Send(internalErrorMessage)
	}
	address := splitData[1]
	if err = dependencies.AccountService.SwitchAddress(
		//nolint:gosec //address index is pretty small number and is never negative
		context.Background(), tgCtx.Sender().ID, uint32(addressIndex)); err != nil {
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Your current address is: %s", address)
	return tgCtx.Send(message)
}

func GetBalance(tgCtx telebot.Context, dependencies *BotServices) error {
	ethBalance, err := dependencies.AccountService.GetBalance(context.Background(), tgCtx.Sender().ID)
	if err != nil {
		return tgCtx.Send(internalErrorMessage)
	}
	message := fmt.Sprintf("Balance for your current address is: %s ETH", ethBalance)
	return tgCtx.Send(message)
}
