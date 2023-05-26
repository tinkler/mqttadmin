package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type loggerLevel uint8

var (
	ConsoleLevel           loggerLevel = LL_SILENT
	LoggerFileLevel        loggerLevel = LL_NORMAL
	ConsoleMessageSeparate             = "\n<<<<<<\n"
)

const (
	_         loggerLevel = iota
	LL_SILENT             // print LML_ERR LML_TITLE message
	LL_NORMAL             // print LML_ERR LML_TITLE LML_WARN message
	LL_LOG                // print LML_ERR LML_TITLE LML_WARN LML_INFO message
	LL_DEBUG              // print LML_ERR LML_TITLE LML_WARN LML_INFO LML_DEBUG message
)

type loggerMessageLevel uint8

const (
	_ loggerMessageLevel = iota
	LML_ERR
	LML_WARN
	LML_INFO
	LML_DEBUG
	LML_TITLE
)

type ColorType uint8

const (
	_ ColorType = iota + 29
	CT_BLACK
	CT_RED
	CT_GREEN
	CT_YELLOW
	CT_BLUE
	CT_PURPLE
	CT_CYAN
	CT_WHITE
)

var sourceDir, errSourceDir, statusSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	sourceDir = regexp.MustCompile(`logger\.go`).ReplaceAllString(file, "")
	errSourceDir, _ = filepath.Abs(sourceDir + "/../../errz")
	errSourceDir = strings.ReplaceAll(errSourceDir, "\\", "/")
	statusSourceDir, _ = filepath.Abs(sourceDir + "/../status")
	statusSourceDir = strings.ReplaceAll(statusSourceDir, "\\", "/")
}

func writeRuntimeMessage() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, sourceDir) && !strings.HasPrefix(file, errSourceDir) && !strings.HasPrefix(file, statusSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return fmt.Sprintf("\033[%d;1m%s:%d\033[0m", CT_BLACK, file, line)
		}
	}
	return ""
}

func wrapMessage(message string, level loggerMessageLevel) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	switch level {
	case LML_TITLE:
		message = fmt.Sprintf("%s \033[%d;1m[标题]\033[0m %s", ts, CT_CYAN, message)
	case LML_DEBUG:
		message = fmt.Sprintf("%s \033[%d;1m[调试]\033[0m %s", ts, CT_GREEN, message)
	case LML_INFO:
		message = fmt.Sprintf("%s \033[%d;1m[日志]\033[0m %s", ts, CT_BLUE, message)
	case LML_WARN:
		message = fmt.Sprintf("%s \033[%d;1m[警告]\033[0m %s", ts, CT_YELLOW, message)
	case LML_ERR:
		message = fmt.Sprintf("%s \033[%d;1m[错误]\033[0m %s", ts, CT_RED, message)
	}
	return message
}

func Format(f interface{}, vs ...any) string {
	var message string
	switch fv := f.(type) {
	case string:
		message = fv
		if len(vs) == 0 {
			return message
		}
		if !strings.Contains(message, "%") {
			message += strings.Repeat(" %v", len(vs))
		}
	default:
		message = fmt.Sprint(f)
		if len(vs) == 0 {
			return message
		}
		message += strings.Repeat(" %v", len(vs))
	}
	if len(vs) > 0 {
		return fmt.Sprintf(message, vs...)
	}
	return fmt.Sprint(message)
}

func Error(f interface{}, vs ...any) {
	b := bytes.NewBufferString(writeRuntimeMessage())
	b.WriteByte('\n')
	b.WriteString(wrapMessage(Format(f, vs...), LML_ERR))
	b.WriteString(ConsoleMessageSeparate)
	os.Stdout.Write(b.Bytes())
	writeToFile(b.Bytes(), LML_ERR)
}

func Warn(f interface{}, vs ...any) {
	b := bytes.NewBufferString(writeRuntimeMessage())
	b.WriteByte('\n')
	b.WriteString(wrapMessage(Format(f, vs...), LML_WARN))
	b.WriteString(ConsoleMessageSeparate)
	if ConsoleLevel > LL_SILENT {
		os.Stdout.Write(b.Bytes())
	}
	if LoggerFileLevel > LL_SILENT {
		writeToFile(b.Bytes(), LML_WARN)
	}
}

func Info(f interface{}, vs ...any) {
	b := bytes.NewBufferString(writeRuntimeMessage())
	b.WriteByte('\n')
	b.WriteString(wrapMessage(Format(f, vs...), LML_INFO))
	b.WriteString(ConsoleMessageSeparate)
	if ConsoleLevel > LL_NORMAL {
		os.Stdout.Write(b.Bytes())
	}
	if LoggerFileLevel > LL_NORMAL {
		writeToFile(b.Bytes(), LML_INFO)
	}
}

func Debug(f interface{}, vs ...any) {
	b := bytes.NewBufferString(writeRuntimeMessage())
	b.WriteByte('\n')
	b.WriteString(wrapMessage(Format(f, vs...), LML_DEBUG))
	b.WriteString(ConsoleMessageSeparate)
	if ConsoleLevel > LL_LOG {
		os.Stdout.Write(b.Bytes())
	}
	if LoggerFileLevel > LL_LOG {
		writeToFile(b.Bytes(), LML_DEBUG)
	}
}

func Title(f interface{}, vs ...any) {
	b := bytes.NewBufferString(wrapMessage(Format(f, vs...), LML_TITLE))
	b.WriteByte('\n')
	os.Stdout.Write(b.Bytes())
	writeToFile(b.Bytes(), LML_TITLE)
}
