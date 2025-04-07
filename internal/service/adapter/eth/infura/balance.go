package infura

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"

	"github.com/pkg/errors"
)

const (
	methodGetBalance = "eth_getBalance"
	blockParamLatest = "latest"
)

var (
	ErrInvalidResultType       = errors.New("invalid result type")
	ErrFailedToParseHexBalance = errors.New("failed to parse hex balance")
)

func (c *Client) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	ctx, span := c.tracer.Start(ctx, "GetBalance")
	defer span.End()
	params := []string{address, blockParamLatest}
	jsonRPCReq := NewRequest(methodGetBalance, params)
	b, err := json.Marshal(jsonRPCReq)
	if err != nil {
		err = errors.Wrap(err, "failed to marshal request")
		span.RecordError(err)
		return big.NewInt(0), err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewBuffer(b))
	if err != nil {
		err = errors.Wrap(err, "failed to create request")
		span.RecordError(err)
		return big.NewInt(0), err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to execute request")
		span.RecordError(err)
		return big.NewInt(0), err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.Errorf("failed to get balance with http status code: %d", resp.StatusCode)
		span.RecordError(err)
		return big.NewInt(0), err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body")
		span.RecordError(err)
		return big.NewInt(0), err
	}
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		err = errors.Wrap(err, "failed to unmarshal response body")
		span.RecordError(err)
		return big.NewInt(0), err
	}
	if response.HasError() {
		err = errors.Errorf("failed to get balance: %s", response.Error.Message)
		span.RecordError(err)
		return big.NewInt(0), err
	}
	balanceStr, ok := response.Result.(string)
	if !ok {
		err = ErrInvalidResultType
		span.RecordError(err)
		return big.NewInt(0), err
	}
	weiBalance, err := parseBalance(balanceStr)
	if err != nil {
		err = errors.Wrap(err, "failed to convert wei to eth")
		span.RecordError(err)
		return big.NewInt(0), err
	}
	return weiBalance, nil
}

func parseBalance(weiStr string) (*big.Int, error) {
	wei := new(big.Int)
	_, ok := wei.SetString(weiStr, 0)
	if !ok {
		return nil, ErrFailedToParseHexBalance
	}
	return wei, nil
}
