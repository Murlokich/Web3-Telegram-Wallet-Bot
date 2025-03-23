package infura

import (
	"Web3-Telegram-Wallet-Bot/internal/config"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
	tracer     trace.Tracer
}

func New(cfg *config.Infura) *Client {
	return &Client{httpClient: http.DefaultClient, endpoint: cfg.Endpoint}
}
