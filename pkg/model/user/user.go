package user

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"google.golang.org/grpc/status"
)

type User struct {
	ID       int64
	Username string
	Email    string
	Profile  *UserProfile `json:",omitempty"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Save(ctx context.Context) error {
	return db.GetDB(ctx).Save(u).Error
}

func (u *User) AddRole(ctx context.Context, role *role.Role) error {
	if role == nil || role.ID == 0 {
		return status.New(400, "role is nil or role id is 0").Err()
	}
	return nil
}
