package user

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"github.com/tinkler/mqttadmin/pkg/status"
)

// User is the user model
type User struct {
	ID       string // ID is the primary key
	Username string
	Email    string
	Profiles []*UserProfile `json:",omitempty" gorm:"-"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "users"
}

// Save saves the user to the database
func (u *User) Save(ctx context.Context) error {
	return db.GetDB(ctx).Save(u).Error
}

// AddRole adds a role to the user
func (u *User) AddRole(ctx context.Context, role *role.Role) error {
	if role == nil || role.ID == 0 {
		return status.New(400, "role is nil or role id is 0")
	}
	return nil
}
