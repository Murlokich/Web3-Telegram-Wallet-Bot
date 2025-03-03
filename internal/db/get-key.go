package db

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func GetChangeLevelKey(ctx context.Context, conn *pgx.Conn, userID int64) (*encryption.EncryptedEntry, error) {
	query := "SELECT change_level_key, nonce FROM user_change_level_key WHERE user_id = $1"
	var entry encryption.EncryptedEntry
	if err := conn.QueryRow(ctx, query, userID).Scan(&entry.Ciphertext, &entry.Nonce); err != nil {
		return nil, errors.Wrap(err, "failed to query row")
	}
	return &entry, nil
}
