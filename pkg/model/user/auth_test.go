package user

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tinkler/mqttadmin/pkg/status"
)

func initEnv() {
	root, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(root, "../../../.env"))
	if err != nil {
		log.Fatal("Error loading .env file:" + err.Error())
	}
}

func TestAuthSignup(t *testing.T) {
	initEnv()
	ctx := context.Background()
	a := &Auth{
		Username: "admin2",
		Password: "admin2",
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
	initEnv()
	ctx := context.Background()
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
