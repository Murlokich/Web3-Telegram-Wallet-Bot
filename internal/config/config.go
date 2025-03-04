package config

type Config struct {
	TelegramBotConfig TelegramBotConfig
	DBConfig          DBConfig
	Encryption        Encryption
}

type TelegramBotConfig struct {
	Token   string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	Timeout int    `envconfig:"TELEGRAM_BOT_POLLING_TIMEOUT" default:"30"`
}

type DBConfig struct {
	URL              string `envconfig:"DB_URL" required:"true"`
	MigrationVersion uint   `envconfig:"DB_MIGRATION_VERSION" required:"true"`
}

type Encryption struct {
	MasterKey string `envconfig:"ENCRYPTION_MASTER_KEY" required:"true"`
}

type Infura struct {
	Endpoint string `envconfig:"INFURA_ENDPOINT" required:"true"`
}
