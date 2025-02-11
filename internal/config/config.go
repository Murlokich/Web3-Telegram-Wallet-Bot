package config

type Config struct {
	TelegramBotConfig TelegramBotConfig
}

type TelegramBotConfig struct {
	Token   string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	Timeout int    `envconfig:"TELEGRAM_BOT_POLLING_TIMEOUT" default:"30"`
}
