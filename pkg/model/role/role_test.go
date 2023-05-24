package role

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func initEnv() {
	root, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(root, "../../../.env"))
	if err != nil {
		log.Fatal("Error loading .env file:" + err.Error())
	}
}

func TestSaveRole(t *testing.T) {
	initEnv()
	role := &Role{
		Name: "user",
	}
	err := role.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(role.ID)
}
