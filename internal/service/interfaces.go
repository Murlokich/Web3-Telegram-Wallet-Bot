package service

import "context"

type AccountService interface {
	CreateAccount(ctx context.Context, userID int64) (string, string, error)
	MigrateAccount(ctx context.Context, mnemonic string, userID int64) (string, error)
	AddNewAddress(ctx context.Context, userID int64) (string, error)
	GetBalance(ctx context.Context, address string) (string, error)
}
