package user

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
)

type UserRole struct {
	ID   uint64
	User *User
	Role *role.Role
}

func (ur *UserRole) Save(ctx context.Context) error {
	return db.GetDB(ctx).Save(ur).Error
}
