package db

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
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

func (c *Client) InsertWallet(ctx context.Context, userID int64,
	mkEntry *encryption.EncryptedEntry, clkEntry *encryption.EncryptedEntry) error {
	query := "INSERT INTO user_wallet VALUES ($1, $2, $3, $4, $5, DEFAULT)"
	tag, err := c.conn.Exec(ctx, query, userID, mkEntry.Ciphertext, mkEntry.Nonce, clkEntry.Ciphertext, clkEntry.Nonce)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	if !tag.Insert() || tag.RowsAffected() != insertWalletAffectedRows {
		return &InsertAnomalyError{Insert: tag.Insert(), RowsAffected: tag.RowsAffected()}
	}
	return nil
}
