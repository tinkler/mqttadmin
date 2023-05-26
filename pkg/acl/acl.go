package acl

import (
	"os"
	"sync"
)

var (
	aclmInstance     *ACLM
	aclmInstanceOnce sync.Once
)

func Aclm() *ACLM {
	aclmInstanceOnce.Do(func() {
		ev := os.Getenv("ACL_MODE")
		if ev == "" {
			aclmInstance = newACLM(AM_CLIENT)
		} else {
			aclmInstance = newACLM(aclmMode(ev))
		}
	})
	return aclmInstance
}
