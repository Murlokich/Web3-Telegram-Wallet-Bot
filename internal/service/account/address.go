package account

import (
	"Web3-Telegram-Wallet-Bot/internal/service"
	"context"
)

const (
	addNewAddressSpanName = "AddNewAddress"
)

func (s *Service) AddNewAddress(ctx context.Context, userID int64) (string, error) {
	ctx, span := s.tracer.Start(ctx, addNewAddressSpanName)
	defer span.End()
	s.Logger.Info("Creating Address")
	res, err := s.DB.AddNewAddress(ctx, userID)
	if err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to add new address: ", err)
		return "", service.ErrInternal
	}
	address, err := s.HDWallet.GetAddress(ctx, res.ChangeLevelKey, res.LastAddressIndex)
	if err != nil {
		span.RecordError(err)
		s.Logger.Error("Failed to get address: ", err)
		return "", service.ErrInternal
	}
	return address, nil
}
