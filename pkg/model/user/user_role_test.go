package user

import (
	"context"
	"testing"

	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
)

func TestSaveUserRole(t *testing.T) {
	initEnv()
	u := &User{
		Username: "admin",
	}
	se := db.DB().Take(u)
	if se.Error != nil {
		t.Fatal(se.Error)
	}
	r := &role.Role{
		Name: "admin",
	}
	se = db.DB().Take(r)
	if se.Error != nil {
		t.Fatal(se.Error)
	}
	userRole := &UserRole{
		User: u,
		Role: r,
	}
	err := userRole.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userRole.ID)
}
