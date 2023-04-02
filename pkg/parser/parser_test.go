package parser

import (
	"testing"
)

func TestParsePackage(t *testing.T) {
	pkg, err := ParsePackage("../model/page")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pkg)
}
