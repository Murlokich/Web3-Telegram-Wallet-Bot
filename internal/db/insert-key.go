package db

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type InsertAnomalyError struct {
	Insert       bool
	RowsAffected int64
}

func (e *InsertAnomalyError) Error() string {
	return fmt.Sprintf("insertion anomaly: expected insert operation true (1 row), got insert %t (%d rows)",
		e.Insert, e.RowsAffected)
}

func InsertKeys(ctx context.Context, conn *pgx.Conn, userID int64,
	mkEntry *encryption.EncryptedEntry, clkEntry *encryption.EncryptedEntry) (retError error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if retError != nil && rollbackErr != nil {
			retError = errors.Wrap(retError, rollbackErr.Error())
		}
	}()

	if err = insertMasterKey(ctx, tx, userID, mkEntry); err != nil {
		return errors.Wrap(err, "failed to insert master key")
	}
	if err = insertChangeLevelKey(ctx, tx, userID, clkEntry); err != nil {
		return errors.Wrap(err, "failed to insert change level key")
	}
	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}

func insertMasterKey(ctx context.Context, tx pgx.Tx, userID int64, entry *encryption.EncryptedEntry) error {
	query := "INSERT INTO user_master_key (user_id, master_key, nonce) VALUES ($1, $2, $3)"
	tag, err := tx.Exec(ctx, query, userID, entry.Ciphertext, entry.Nonce)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	if !tag.Insert() || tag.RowsAffected() != 1 {
		return &InsertAnomalyError{Insert: tag.Insert(), RowsAffected: tag.RowsAffected()}
	}
	return nil
}

func insertChangeLevelKey(ctx context.Context, tx pgx.Tx, userID int64, entry *encryption.EncryptedEntry) error {
	query := "INSERT INTO user_change_level_key (user_id, change_level_key, nonce) VALUES ($1, $2, $3)"
	tag, err := tx.Exec(ctx, query, userID, entry.Ciphertext, entry.Nonce)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	if !tag.Insert() || tag.RowsAffected() != 1 {
		return &InsertAnomalyError{Insert: tag.Insert(), RowsAffected: tag.RowsAffected()}
	}
	return nil
}
