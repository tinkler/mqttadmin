package parser

import (
	"os"
	"testing"
)

func TestGetModulePath(t *testing.T) {
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	if GetModulePath() != "github.com/tinkler/mqttadmin" {
		t.Fail()
	}
}
