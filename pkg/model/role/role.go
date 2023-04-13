package role

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/db"
)

type Role struct {
	ID   int64
	Name string
}

func (r *Role) SaveRole(ctx context.Context) error {
	return db.GetDB(ctx).Save(r).Error
}
