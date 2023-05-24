package user

import (
	"context"
	"testing"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
)

func TestUserAddRole(t *testing.T) {
	initEnv()
	u := &User{
		Username: "admin",
	}
	err := u.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	r := &role.Role{
		Name: "user",
	}
	se := db.DB().Take(r)
	if se.Error != nil {
		t.Fatal(se.Error)
	}
	err = u.AddRole(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGet(t *testing.T) {
	initEnv()
	u := &User{
		Username: "admin",
	}
	err := u.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestUserGetRoles(t *testing.T) {
	initEnv()
	u := &User{
		Username: "admin",
	}
	err := u.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	err = u.GetRoles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range u.Roles {
		t.Log(r)
	}
}
