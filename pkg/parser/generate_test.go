package parser

import (
	"os"
	"testing"
)

func TestGenerateGoCode(t *testing.T) {
	modulePath := GetModulePath() + "/pkg"
	pkg, err := ParsePackage("../model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("../model/role", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("../model/page", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}
	err = GenerateChiRoutes("../route", pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateChiRoutes("../route", pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateChiRoutes("../route", pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateProtoFile(t *testing.T) {
	_ = os.Chdir("../../")
	modulePath := GetModulePath()
	pkg, err := ParsePackage("./pkg/model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("./pkg/model/role", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("./pkg/model/page", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}

	basePath := modulePath
	err = GenerateProtoFile("./api/proto", basePath, pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateProtoFile("./api/proto", basePath, pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateProtoFile("./api/proto", basePath, pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateGsrv(t *testing.T) {
	_ = os.Chdir("../../")
	modulePath := GetModulePath()
	pkg, err := ParsePackage("./pkg/model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("./pkg/model/role", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("./pkg/model/page", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}

	err = GenerateGsrv("./pkg/gsrv", modulePath, pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateGsrv("./pkg/gsrv", modulePath, pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateGsrv("./pkg/gsrv", modulePath, pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateTSCode(t *testing.T) {
	modulePath := GetModulePath() + "/pkg"
	pkg, err := ParsePackage("../model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("../model/role", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("../model/page", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}
	err = GenerateTSCode("../../static/ts", pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	// dependency test
	err = GenerateTSCode("../../static/ts", pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	// array test
	err = GenerateTSCode("../../static/ts", pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGenerateTSAngularDelonCode(t *testing.T) {
	modulePath := GetModulePath() + "/pkg"
	pkg, err := ParsePackage("../model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("../model/role", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("../model/page", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}
	err = GenerateTSAngularDelonCode("../../static/angular_delon/mqtt", pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	// dependency test
	err = GenerateTSAngularDelonCode("../../static/angular_delon/mqtt", pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	// array test
	err = GenerateTSAngularDelonCode("../../static/angular_delon/mqtt", pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGenerateDartCode(t *testing.T) {
	modulePath := GetModulePath() + "/pkg"
	pkg, err := ParsePackage("../model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg2, err := ParsePackage("../model/role", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := ParsePackage("../model/page", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	pkgs := map[string]*Package{"role": pkg2, "user": pkg, "page": pkg3}
	err = GenerateDartCode("../../static/dart", pkg, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateDartCode("../../static/dart", pkg2, pkgs)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateDartCode("../../static/dart", pkg3, pkgs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateSwiftCode(t *testing.T) {
	modulePath := GetModulePath() + "/pkg"
	pkg, err := ParsePackage("../model/user", modulePath)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateSwiftCode("../../", pkg)
	if err != nil {
		t.Fatal(err)
	}
}
