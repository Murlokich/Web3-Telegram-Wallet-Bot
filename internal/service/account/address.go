package account

import (
	"Web3-Telegram-Wallet-Bot/internal/service"
	"context"

	"github.com/pkg/errors"
)

func (s *Service) AddNewAddress(ctx context.Context, userID int64) (string, error) {
	ctx, span := s.tracer.Start(ctx, "AddNewAddress")
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

func (s *Service) GetAllAddresses(ctx context.Context, userID int64) ([]string, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllAddresses")
	defer span.End()
	s.Logger.Info("Getting all addresses")
	data, err := s.DB.GetChangeLevelKey(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "failed to get change level key from db")
		span.RecordError(err)
		return nil, err
	}
	numOfAddresses := data.LastAddressIndex + 1
	addresses := make([]string, numOfAddresses)
	for i := range numOfAddresses {
		addresses[i], err = s.HDWallet.GetAddress(ctx, data.ChangeLevelKey, i)
		if err != nil {
			err = errors.Wrap(err, "failed to get address from HDWallet")
			span.RecordError(err)
			return nil, err
		}
	}
	return addresses, nil
}

func (s *Service) SwitchAddress(ctx context.Context, userID int64, addressIndex uint32) error {
	ctx, span := s.tracer.Start(ctx, "SwitchAddress")
	defer span.End()
	if err := s.DB.UpdateCurrentAddress(ctx, userID, addressIndex); err != nil {
		err = errors.Wrap(err, "failed to update current address")
		span.RecordError(err)
		return err
	}
	return nil
}
