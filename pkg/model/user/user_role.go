package user

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"github.com/tinkler/mqttadmin/pkg/status"
)

type UserRole struct {
	ID   string
	User *User      `gorm:"foreignkey:UserID"`
	Role *role.Role `gorm:"foreignkey:RoleID"`
}

func (ur *UserRole) TableName() string {
	return "v1.user_role"
}

func (ur *UserRole) Save(ctx context.Context) error {
	if ur.User == nil || ur.User.ID == "" {
		return status.New(400, "user is nil or user id is 0")
	}
	return db.DB().Select("user_id", "role_id").Save(ur).Error
}
