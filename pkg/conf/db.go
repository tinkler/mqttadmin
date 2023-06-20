package conf

import (
	"os"
	"sync"
)

type DbConfig struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

var (
	dbOnce sync.Once
	dbIns  *DbConfig
)

func newDbConfig() *DbConfig {
	dbOnce.Do(func() {
		// Dsn database connection string:
		driver := os.Getenv("DB_DRIVER")
		dsn := os.Getenv("DB_DSN")
		dbIns = &DbConfig{
			Driver: driver,
			Dsn:    dsn,
		}

	})

	return dbIns
}
