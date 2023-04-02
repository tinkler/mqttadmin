package user

import (
	"context"
	"fmt"

	"github.com/tinkler/mqttadmin/pkg/model/role"
)

type User struct {
	Ctx     context.Context `json:"-"`
	ID      int64
	Name    string
	Email   string
	Age     int
	Profile *UserProfile `json:",omitempty"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Save(ctx context.Context) error {
	u.ID = 10
	return nil
}

func (u *User) AddRole(ctx context.Context, role *role.Role) error {
	fmt.Println(role)
	return nil
}
