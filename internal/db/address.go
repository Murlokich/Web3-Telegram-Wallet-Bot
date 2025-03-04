package db

import (
	"Web3-Telegram-Wallet-Bot/internal/encryption"
	"context"

	"github.com/pkg/errors"
)

func (c *Client) AddNewAddress(ctx context.Context, userID int64) (*encryption.EncryptedEntry, uint32, error) {
	query := `
		UPDATE user_wallet SET last_address_index = last_address_index + 1
		WHERE user_id = $1 RETURNING change_level_key, clk_nonce, last_address_index
    `
	var entry encryption.EncryptedEntry
	var lastAddressIndex uint32
	if err := c.conn.QueryRow(ctx, query, userID).Scan(&entry.Ciphertext, &entry.Nonce, &lastAddressIndex); err != nil {
		return nil, 0, errors.Wrap(err, "failed to query row")
	}
	return &entry, lastAddressIndex, nil
}

func (c *Client) GetChangeLevelKey(ctx context.Context, userID int64) (*encryption.EncryptedEntry, uint32, error) {
	query := `
		SELECT change_level_key, clk_nonce, last_address_index 
		FROM user_wallet
		WHERE user_id = $1`
	var entry encryption.EncryptedEntry
	var lastAddressIndex uint32
	if err := c.conn.QueryRow(ctx, query, userID).Scan(&entry.Ciphertext, &entry.Nonce, &lastAddressIndex); err != nil {
		return nil, 0, errors.Wrap(err, "failed to query row")
	}
	return &entry, lastAddressIndex, nil
}
