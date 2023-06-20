package db

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/conf"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type contextKey string

const (
	dbKey contextKey = "db"
)

func NewDB(cnf *conf.Conf, gormCnf *gorm.Config) (*gorm.DB, error) {
	switch cnf.Db.Driver {
	case "mssql":
		return gorm.Open(sqlserver.Open(cnf.Db.Dsn), gormCnf)
	case "postgres":
		fallthrough
	default:
		return gorm.Open(postgres.Open(cnf.Db.Dsn), gormCnf)
	}
}

func GetDB(ctx context.Context) *gorm.DB {
	i := ctx.Value(dbKey)
	if i == nil {
		panic("db not found in context")
	}
	return i.(*gorm.DB)
}
