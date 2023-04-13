package permission

import (
	"context"
	"net/http"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"github.com/tinkler/mqttadmin/pkg/status"
)

type RolePermission struct {
	ID         int64
	Role       *role.Role
	Permission *Permission
}

func (rp *RolePermission) SavePermission(ctx context.Context) error {
	if rp.Permission == nil || rp.Permission.ID == 0 {
		return status.New(http.StatusBadRequest, "permission is nil or permission id is 0")
	}
	return db.GetDB(ctx).Save(rp).Error
}

func (rp *RolePermission) DeletePermission(ctx context.Context) error {
	if rp.ID == 0 {
		return status.New(http.StatusBadRequest, "role permission id is 0")
	}
	return db.GetDB(ctx).Delete(rp).Error
}
