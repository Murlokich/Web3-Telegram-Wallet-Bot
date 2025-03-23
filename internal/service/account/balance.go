package account

import (
	"context"
	"math/big"

	"github.com/pkg/errors"
)

const (
	weiConversionBase  = 10
	weiConversionPower = 18
	ethFloatPrecision  = 8
)

func (s *Service) GetBalance(ctx context.Context, address string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "GetBalance")
	defer span.End()

	weiBalance, err := s.ETH.GetBalance(ctx, address)
	if err != nil {
		err = errors.Wrap(err, "failed to get balance for the address")
		span.RecordError(err)
		return "", err
	}
	return convertWeiToETH(weiBalance), nil
}

func convertWeiToETH(weiBalance *big.Int) string {
	weiFloat := new(big.Float).SetInt(weiBalance)
	// divisor for conversion Wei to ETH is 10**18
	ethDivisor := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(weiConversionBase), big.NewInt(weiConversionPower), nil),
	)
	eth := new(big.Float).Quo(weiFloat, ethDivisor)
	return eth.Text('f', ethFloatPrecision)
}
