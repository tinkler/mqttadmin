package user

import (
	"context"
	"errors"
	"testing"

	"github.com/tinkler/mqttadmin/pkg/conf"
	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/status"
	"gorm.io/gorm"
)

func getWrappedContext() context.Context {
	ctx := context.Background()
	testConf := &conf.Conf{}
	testConf.Db = &conf.DbConfig{
		Dsn: "host=localhost user=clans password=clans4105 dbname=clans port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	}
	dbInst, err := db.NewDB(testConf, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	ctx = db.WithValue(ctx, dbInst)
	return ctx
}

func TestAuthSignup(t *testing.T) {
	ctx := getWrappedContext()
	a := &Auth{
		Username: "admin2",
		Password: "admin",
	}
	u, err := a.Signup(ctx)
	if err != nil && !errors.Is(err, status.Ok("username already exists")) {
		t.Fatal(err)
	}
	if errors.Is(err, status.Ok(ErrMsgNameAreadyExist)) {
		return
	}
	if u.ID == "" {
		t.Fatalf("user id is empty")
	}
}

func TestAuthSignin(t *testing.T) {
	ctx := getWrappedContext()
	a := &Auth{
		Username: "admin2",
		Password: "admin",
	}
	u, err := a.Signin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == "" {
		t.Fatalf("user id is empty")
	}
}
