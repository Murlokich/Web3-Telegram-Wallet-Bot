package postgres

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Client struct {
	conn *pgx.Conn
}

func New(ctx context.Context, config *config.DBConfig) (*Client, error) {
	conn, err := pgx.Connect(ctx, config.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close(context.Background())
}
