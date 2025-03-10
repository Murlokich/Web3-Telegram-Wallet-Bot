package account

import (
	"Web3-Telegram-Wallet-Bot/internal/service"
	"context"
)

func (s *Service) AddNewAddress(ctx context.Context, userID int64) (string, error) {
	s.Logger.Info("Creating Address")
	res, err := s.DB.AddNewAddress(ctx, userID)
	if err != nil {
		s.Logger.Error("Failed to add new address: ", err)
		return "", service.ErrInternal
	}
	address, err := s.HDWallet.GetAddress(res.ChangeLevelKey, res.LastAddressIndex)
	if err != nil {
		s.Logger.Error("Failed to get address: ", err)
		return "", service.ErrInternal
	}
	return address, nil
}
