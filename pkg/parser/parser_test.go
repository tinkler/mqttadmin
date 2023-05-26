package parser

import (
	"os"
	"testing"
)

func TestParsePackage(t *testing.T) {
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	pkg, err := ParsePackage("pkg/model/user", GetModulePath())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pkg)
}
