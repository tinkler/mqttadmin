package user

import (
	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"github.com/tinkler/mqttadmin/pkg/status"
)

func addRole(u *User, role *role.Role) error {
	if u.ID == "" {
		return status.New(400, "user id is 0")
	}
	if role == nil || (role.ID == "" && role.Name == "") {
		return status.New(400, "role is nil or role id is empty")
	}
	if role.ID == "" {
		if err := db.DB().First(role).Error; err != nil {
			return err
		}
	}
	return db.DB().Model(u).Association("Roles").Append(role)
}

func removeRole(u *User, role *role.Role) error {
	if u.ID == "" {
		return status.New(400, "user id is 0")
	}
	if role == nil || (role.ID == "" && role.Name == "") {
		return status.New(400, "role is nil or role id is empty")
	}
	if role.ID == "" {
		if err := db.DB().First(role).Error; err != nil {
			return err
		}
	}
	return db.DB().Model(u).Association("Roles").Delete(role)
}

func (u *User) GetRolesStrings() []string {
	role := make([]string, len(u.Roles))
	for i, v := range u.Roles {
		role[i] = v.Name
	}
	return role
}
