package user

import (
	"context"
	"net/http"

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
	Roles    []*role.Role   `gorm:"many2many:v1.user_role;"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "v1.user"
}

// Save saves the user to the database
func (u *User) Save(ctx context.Context) error {
	if !HasRole(ctx, RoleAdmin) {
		return status.New(http.StatusForbidden, "no permission")
	}
	return db.DB().Save(u).Error
}

// AddRole adds a role to the user
func (u *User) AddRole(ctx context.Context, role *role.Role) error {
	if u.ID == "" {
		return status.New(400, "user id is 0")
	}
	if role == nil || role.ID == "" {
		return status.New(400, "role is nil or role id is empty")
	}
	if !HasRole(ctx, RoleAdmin) {
		return status.StatusForbidden
	}
	return db.DB().Model(u).Association("Roles").Append(role)
}

// Get gets the user from the database
// Only admin can get other user's information
func (u *User) Get(ctx context.Context) error {
	if HasRole(ctx, RoleAdmin) {
		if u.ID == "" {
			userID := GetUserID(ctx)
			if userID == "" {
				return status.New(400, "user id is 0")
			} else {
				u.ID = userID
			}
		}
		return db.DB().First(u).Error
	}
	if HasRole(ctx, RoleUser) {
		u.ID = GetUserID(ctx)
		if u.ID == "" {
			return status.StatusInternalServer
		}
		return db.DB().Where("id = ?", u.ID).First(u).Error
	}
	return status.StatusUnauthorized
}

// GetRoles gets the roles of the user
func (u *User) GetRoles(ctx context.Context) error {
	if u.ID == "" {
		return status.New(400, "user id is 0")
	}
	u.Roles = []*role.Role{}
	if err := db.DB().Model(u).Association("Roles").Find(&u.Roles); err != nil {
		return err
	}
	return nil
}
