package postgres

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Client struct {
	conn   *pgx.Conn
	tracer trace.Tracer
}

func New(ctx context.Context, tracer trace.Tracer, config *config.DBConfig) (*Client, error) {
	conn, err := pgx.Connect(ctx, config.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}
	return &Client{conn: conn, tracer: tracer}, nil
}

func (c *Client) Close() error {
	return c.conn.Close(context.Background())
}
