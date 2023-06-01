package role

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/db"
	"gorm.io/gorm"
)

type Role struct {
	ID       string
	Name     string
	DeleteAt gorm.DeletedAt
}

func (r *Role) TableName() string {
	return "authv1.role"
}

func (r *Role) Save(ctx context.Context) error {
	return db.DB().Select("name").Save(r).Error
}
