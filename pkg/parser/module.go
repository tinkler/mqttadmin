package parser

import (
	"os/exec"
	"strings"
)

func GetModulePath() string {
	output, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(output))
}
