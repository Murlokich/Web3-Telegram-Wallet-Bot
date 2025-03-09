package postgres

import (
	"Web3-Telegram-Wallet-Bot/internal/repository"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

const (
	insertWalletAffectedRows = 1
)

type InsertAnomalyError struct {
	Insert       bool
	RowsAffected int64
}

func (e *InsertAnomalyError) Error() string {
	return fmt.Sprintf("insertion anomaly: expected insert operation true (%d row), got insert %t (%d rows)",
		insertWalletAffectedRows, e.Insert, e.RowsAffected)
}

func (c *Client) InsertWallet(ctx context.Context, record *repository.WalletEncryptedRecord) error {
	query := "INSERT INTO user_wallet VALUES ($1, $2, $3, $4, $5, $6)"
	tag, err := c.conn.Exec(ctx, query, record.UserID, record.MasterKey.Ciphertext, record.MasterKey.Nonce,
		record.ChangeLevelKey.Ciphertext, record.ChangeLevelKey.Nonce, record.LastAddressIndex)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	if !tag.Insert() || tag.RowsAffected() != insertWalletAffectedRows {
		return &InsertAnomalyError{Insert: tag.Insert(), RowsAffected: tag.RowsAffected()}
	}
	return nil
}
