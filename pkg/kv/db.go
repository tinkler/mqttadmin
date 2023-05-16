package kv

import (
	"database/sql"
	"time"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/logger"
)

type dbBackend struct {
}

func (b dbBackend) Get(key string) (string, bool) {
	var value sql.NullString
	if err := db.DB().Raw(`SELECT value FROM sys.server_cache WHERE key = $1
		AND expire_time >= now() AND deleted = false
		UNION ALL 
	SELECT value FROM sys.server_cache WHERE key = $1
		AND expire_time IS NULL AND deleted = false
	`, key).Find(&value).Error; err != nil {
		logger.Error(err)
	}
	return value.String, value.Valid
}

func (b dbBackend) Set(key string, value string, expireTime *time.Time) {
	if expireTime == nil {
		if err := db.DB().Exec(`INSERT INTO sys.server_cache(key,value,update_time)
		VALUES($1,$2,$3) ON CONFLICT(key)
		DO
			UPDATE SET value = $2
			,expire_time = null
			,update_time = $3
		`, key, value, time.Now()).Error; err != nil {
			logger.Error(err)
		}
		return
	}
	if err := db.DB().Exec(`INSERT INTO sys.server_cache(key,value,expire_time,update_time)
		VALUES($1,$2,$3,$4) ON CONFLICT(key)
	DO
		UPDATE SET value = $2
		,expire_time = $3
		,update_time = $4`, key, value, expireTime,
		time.Now(),
	).Error; err != nil {
		logger.Error(err)
	}
}
