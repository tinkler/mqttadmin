package db

import (
	"github.com/tinkler/mqttadmin/pkg/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cnf *conf.Conf, gormCnf *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cnf.Db.Dsn), gormCnf)
}
