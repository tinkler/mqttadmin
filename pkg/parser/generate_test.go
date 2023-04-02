package parser

import "testing"

func TestGenerateGoCode(t *testing.T) {
	pkg, err := ParsePackage("../model/user")
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("../model/role")
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("../model/page")
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}
	err = GenerateRoutes("../../", pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateRoutes("../../", pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateRoutes("../../", pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateTSCode(t *testing.T) {
	pkg, err := ParsePackage("../model/user")
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("../model/role")
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("../model/page")
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}
	err = GenerateTSCode("../../", pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	// dependency test
	err = GenerateTSCode("../../", pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	// array test
	err = GenerateTSCode("../../", pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGenerateDartCode(t *testing.T) {
	pkg, err := ParsePackage("../model/user")
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateDartCode("../../", pkg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateSwiftCode(t *testing.T) {
	pkg, err := ParsePackage("../model/user")
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateSwiftCode("../../", pkg)
	if err != nil {
		t.Fatal(err)
	}
}
