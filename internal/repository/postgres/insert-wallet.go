package postgres

import (
	"Web3-Telegram-Wallet-Bot/internal/repository"
	"context"

	"github.com/pkg/errors"
)

func (c *Client) InsertWallet(ctx context.Context, record *repository.WalletEncryptedRecord) error {
	ctx, span := c.tracer.Start(ctx, "InsertWallet")
	defer span.End()
	query := "INSERT INTO user_wallet VALUES ($1, $2, $3, $4, $5, $6)"
	tag, err := c.conn.Exec(ctx, query, record.UserID, record.MasterKey.Ciphertext, record.MasterKey.Nonce,
		record.ChangeLevelKey.Ciphertext, record.ChangeLevelKey.Nonce, record.LastAddressIndex)
	if err != nil {
		err = errors.Wrap(err, "failed to execute query")
		span.RecordError(err)
		return err
	}
	if !tag.Insert() || tag.RowsAffected() != insertWalletAffectedRows {
		err = &InsertAnomalyError{
			Insert:               tag.Insert(),
			RowsAffected:         tag.RowsAffected(),
			ExpectedRowsAffected: insertWalletAffectedRows,
		}
		span.RecordError(err)
		return err
	}
	return nil
}
