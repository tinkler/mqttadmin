package conf

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
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
		// Set all variables to the config instance.
		cnfByt, err := os.ReadFile("db.yaml")
		if err != nil {
			panic(err)
		}

		dbIns = &DbConfig{}
		err = yaml.Unmarshal(cnfByt, dbIns)
		if err != nil {
			panic(err)
		}

	})

	return dbIns
}
