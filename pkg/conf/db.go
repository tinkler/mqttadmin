package conf

import (
	"os"
	"sync"
)

type DbConfig struct {
	Dsn string `yaml:"dsn"`
}

var (
	dbOnce sync.Once
	dbIns  *DbConfig
)

func newDbConfig() *DbConfig {
	dbOnce.Do(func() {
		// Dsn database connection string:
		dsn := os.Getenv("DB_DSN")
		dbIns = &DbConfig{
			Dsn: dsn,
		}

	})

	return dbIns
}
