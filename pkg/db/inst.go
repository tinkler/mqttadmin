package db

import (
	"sync"

	"github.com/tinkler/mqttadmin/pkg/conf"
	"gorm.io/gorm"
)

var (
	dbInstance     *gorm.DB
	dbInstanceOnce sync.Once
)

func DB() *gorm.DB {
	dbInstanceOnce.Do(func() {
		var err error
		dbInstance, err = NewDB(conf.NewConf(), conf.NewGormConfig())
		if err != nil {
			panic(err)
		}
	})
	return dbInstance
}
