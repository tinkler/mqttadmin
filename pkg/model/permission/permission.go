package permission

import (
	"context"
	"net/http"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/status"
)

type Permission struct {
	ID   int64
	Name string
}

func (p *Permission) SavePermission(ctx context.Context) error {
	return db.GetDB(ctx).Save(p).Error
}

func (p *Permission) DeletePermission(ctx context.Context) error {
	if p.ID == 0 {
		return status.New(http.StatusBadRequest, "permission id is 0")
	}
	{
		var roleCount int64
		if err := db.GetDB(ctx).Model(&RolePermission{}).Where("permission_id = ?", p.ID).Count(&roleCount).Error; err != nil {
			return err
		}
		if roleCount > 0 {
			return status.New(http.StatusOK, "permission is in use")
		}
	}
	{
		var userCount int64
		if err := db.GetDB(ctx).Model(&UserPermission{}).Where("permission_id = ?", p.ID).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount > 0 {
			return status.New(http.StatusOK, "permission is in use")
		}
	}

	return db.GetDB(ctx).Delete(p).Error
}
