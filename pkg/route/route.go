package route

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var source string

func init() {
	_, file, _, _ := runtime.Caller(0)
	source = regexp.MustCompile(`route\.go`).ReplaceAllString(file, "")
	source, _ = filepath.Abs(source + "/../../")
	source = strings.ReplaceAll(source, "\\", "/")
}

var routePathMap = make(map[string]string)

func GetPathDebugLine(pattern string) map[string]string {
	pattern = strings.TrimSuffix(pattern, "/")
	newMap := make(map[string]string)
	for k, v := range routePathMap {
		newMap[pattern+k] = source + "/" + v
	}
	return newMap
}

type Model[T any, S any] struct {
	Data T `json:",omitempty"`
	Args S `json:",omitempty"`
}

type Res[T any, S any] struct {
	Data T `json:",omitempty"`
	Resp S `json:",omitempty"`
}
