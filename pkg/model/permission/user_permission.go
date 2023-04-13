package permission

import (
	"context"
	"net/http"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/user"
	"github.com/tinkler/mqttadmin/pkg/status"
)

type UserPermission struct {
	ID         int64
	User       *user.User
	Permission *Permission
}

func (up *UserPermission) SavePermission(ctx context.Context) error {
	if up.Permission == nil || up.Permission.ID == 0 {
		return status.New(http.StatusBadRequest, "permission is nil or permission id is 0")
	}
	if up.User == nil || up.User.ID == 0 {
		return status.New(http.StatusBadRequest, "user is nil or user id is 0")
	}
	return db.GetDB(ctx).Save(up).Error
}

func (up *UserPermission) DeletePermission(ctx context.Context) error {
	if up.ID == 0 {
		return status.New(http.StatusBadRequest, "user permission id is 0")
	}
	return db.GetDB(ctx).Delete(up).Error
}
