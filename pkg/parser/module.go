package parser

import (
	"os"
	"os/exec"
	"strings"
)

func GetModulePath() string {
	output, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		panic(err)
	}
	moudles := strings.Split(string(output), "\n")
	root, _ := os.Getwd()
	root = strings.ReplaceAll(root, "\\", "/")
	for _, m := range moudles {
		if strings.HasSuffix(root, m) {
			return m
		}
	}
	return strings.TrimSpace(string(output))
}
