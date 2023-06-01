package user

import (
	"context"
	"net/http"

	"github.com/tinkler/mqttadmin/pkg/acl"
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
	Roles    []*role.Role   `gorm:"many2many:authv1.user_role;"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "authv1.user"
}

// Save saves the user to the database
func (u *User) Save(ctx context.Context) error {
	if !acl.HasRole(ctx, acl.RoleAdmin) {
		return status.New(http.StatusForbidden, "no permission")
	}
	return db.DB().Save(u).Error
}

// AddRole adds a role to the user
func (u *User) AddRole(ctx context.Context, role *role.Role) error {
	if !acl.HasRole(ctx, acl.RoleAdmin) {
		return status.StatusForbidden()
	}
	return addRole(u, role)
}

func (u *User) RemoveRole(ctx context.Context, role *role.Role) error {
	if !acl.HasRole(ctx, acl.RoleAdmin) {
		return status.StatusForbidden()
	}
	return removeRole(u, role)
}

// Get gets the user from the database
// Only admin can get other user's information
func (u *User) Get(ctx context.Context) error {
	if acl.HasRole(ctx, acl.RoleAdmin) {
		if u.ID == "" {
			userID := acl.GetUserID(ctx)
			if userID == "" {
				return status.New(400, "user id is 0")
			} else {
				u.ID = userID
			}
		}
		return db.DB().First(u).Error
	}
	if acl.HasRole(ctx, acl.RoleUser) {
		u.ID = acl.GetUserID(ctx)
		if u.ID == "" {
			return status.StatusInternalServer()
		}
		return db.DB().First(u).Error
	}
	return status.StatusUnauthorized()
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
