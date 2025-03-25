package postgres

import (
	"Web3-Telegram-Wallet-Bot/internal/repository"
	"context"

	"github.com/pkg/errors"
)

func (c *Client) AddNewAddress(ctx context.Context, userID int64) (*repository.AddressManagementEncryptedData, error) {
	ctx, span := c.tracer.Start(ctx, "AddNewAddress")
	defer span.End()
	query := `
		UPDATE user_wallet SET last_address_index = last_address_index + 1
		WHERE user_id = $1 RETURNING change_level_key, clk_nonce, last_address_index
    `
	var res repository.AddressManagementEncryptedData
	if err := c.conn.QueryRow(ctx, query, userID).Scan(
		&res.ChangeLevelKey.Ciphertext,
		&res.ChangeLevelKey.Nonce,
		&res.LastAddressIndex); err != nil {
		err = errors.Wrap(err, "failed to query row")
		span.RecordError(err)
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetChangeLevelKey(
	ctx context.Context, userID int64) (*repository.AddressManagementEncryptedData, error) {
	ctx, span := c.tracer.Start(ctx, "GetChangeLevelKey")
	defer span.End()
	query := `
		SELECT change_level_key, clk_nonce, current_address_index, last_address_index 
		FROM user_wallet
		WHERE user_id = $1`
	var res repository.AddressManagementEncryptedData
	if err := c.conn.QueryRow(ctx, query, userID).Scan(
		&res.ChangeLevelKey.Ciphertext,
		&res.ChangeLevelKey.Nonce,
		&res.CurrentAddressIndex,
		&res.LastAddressIndex); err != nil {
		err = errors.Wrap(err, "failed to query row")
		span.RecordError(err)
		return nil, err
	}
	return &res, nil
}

func (c *Client) UpdateCurrentAddress(ctx context.Context, userID int64, addressIndex uint32) error {
	ctx, span := c.tracer.Start(ctx, "UpdateCurrentAddress")
	defer span.End()
	query := `
		UPDATE user_wallet 
		SET current_address_index = $2
		WHERE user_id = $1`
	tag, err := c.conn.Exec(ctx, query, userID, addressIndex)
	if err != nil {
		err = errors.Wrap(err, "failed to execute query")
		span.RecordError(err)
		return err
	}
	if !tag.Update() || tag.RowsAffected() != updateAddressAffectedRows {
		err = &InsertAnomalyError{
			Insert:               tag.Insert(),
			RowsAffected:         tag.RowsAffected(),
			ExpectedRowsAffected: updateAddressAffectedRows,
		}
		span.RecordError(err)
		return err
	}
	return nil
}
