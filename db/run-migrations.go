package db

import (
	"Web3-Telegram-Wallet-Bot/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // required for migrate.New()
	_ "github.com/golang-migrate/migrate/v4/source/file"       // required for migrate.New()
	"github.com/pkg/errors"
)

func RunMigrations(dbConfig *config.DBConfig) error {
	m, err := migrate.New("file://db/migrations", dbConfig.URL)
	if err != nil {
		return errors.Wrap(err, "failed to create migration")
	}
	err = m.Migrate(dbConfig.MigrationVersion)
	if err != nil {
		return errors.Wrap(err, "failed to migrate database")
	}
	errSrc, errDB := m.Close()
	if errSrc != nil {
		return errors.Wrap(err, "failed to close migration source")
	}
	if errDB != nil {
		return errors.Wrap(err, "failed to close database connection")
	}
	return nil
}
