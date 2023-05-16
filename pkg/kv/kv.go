package kv

import "time"

type kv struct {
	backend Driver
}

func Get(key string) (string, bool) {
	return Kv().backend.Get(key)
}

func Set[T time.Time | *time.Time | int64 | int](key string, value string, expiredTime T) {
	var t interface{} = expiredTime
	var et *time.Time
	switch v := t.(type) {
	case time.Time:
		et = &v
	case *time.Time:
		et = v
	case int64:
		if v != 0 {
			t := time.Unix(v, 0)
			et = &t
		}

	case int:
		if v != 0 {
			t := time.Unix(int64(v), 0)
			et = &t
		}
	}
	Kv().backend.Set(key, value, et)
}
