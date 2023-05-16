package kv

import "time"

type Driver interface {
	Get(key string) (string, bool)
	Set(key string, value string, expireTime *time.Time)
}
