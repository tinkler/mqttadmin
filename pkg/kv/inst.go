package kv

import (
	"sync"
)

var (
	kvInstance     *kv
	kvInstanceOnce sync.Once
)

func Kv() *kv {
	kvInstanceOnce.Do(func() {
		kvInstance = &kv{
			backend: dbBackend{},
		}
	})
	return kvInstance
}
