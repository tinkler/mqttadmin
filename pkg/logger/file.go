package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	currentToFile toFile = TF_NO
)

type toFile uint8

const (
	_ toFile = iota
	TF_NO
	TF_ONE
	TF_SEPARATE
)

var (
	fileLoggerInstance     fileLogger
	fileLoggerInstanceOnce sync.Once
)

type fileLogger interface {
	FileTime() time.Time
	NewFile()
	Log(message []byte, level loggerMessageLevel)
	Close() error
}

type oneFileLogger struct {
	rootPath string
	fileTime time.Time
	file     LogFile
}

func (fl *oneFileLogger) FileTime() time.Time {
	return fl.fileTime
}

func (fl *oneFileLogger) NewFile() {
	newFileTime := time.Now()
	var lf LogFile = &logFile{
		path: filepath.Join(fl.rootPath, newFileTime.Format("2006-01-02")+".log"),
	}
	err := fl.file.Close()
	if err != nil {
		return
	}
	fl.file = lf
	fl.fileTime = newFileTime
}

func (fl *oneFileLogger) Log(message []byte, level loggerMessageLevel) {
	fl.file.Write(message)
}

func (fl *oneFileLogger) Close() error {
	return fl.file.Close()
}

type separateFileLogger struct {
	rootPath                               string
	fileTime                               time.Time
	errFile, warnFile, infoFile, debugFile LogFile
}

func (fl *separateFileLogger) FileTime() time.Time {
	return fl.fileTime
}

func (fl *separateFileLogger) NewFile() {
	newFileTime := time.Now()
	newFilePath := func(name string) string {
		return fmt.Sprintf(filepath.Join(fl.rootPath, newFileTime.Format("2006-01-02")+".%s.log"), name)
	}
	var errFile LogFile = &logFile{
		path: newFilePath("err"),
	}
	err := fl.errFile.Close()
	if err != nil {
		return
	}
	fl.errFile = errFile

	var warnFile LogFile = &logFile{
		path: newFilePath("warn"),
	}
	err = fl.warnFile.Close()
	if err != nil {
		return
	}
	fl.warnFile = warnFile

	var infoFile LogFile = &logFile{
		path: newFilePath("info"),
	}
	err = fl.infoFile.Close()
	if err != nil {
		return
	}
	fl.infoFile = infoFile

	var debugFile LogFile = &logFile{
		path: newFilePath("debug"),
	}
	err = fl.debugFile.Close()
	if err != nil {
		return
	}
	fl.debugFile = debugFile

	fl.fileTime = newFileTime
}

func (fl separateFileLogger) Log(message []byte, level loggerMessageLevel) {
	switch level {
	case LML_ERR:
		fl.errFile.Write(message)
	case LML_WARN:
		fl.warnFile.Write(message)
	case LML_INFO:
		fl.warnFile.Write(message)
	case LML_DEBUG:
		fl.warnFile.Write(message)
	}
}

func (fl separateFileLogger) Close() error {
	err := fl.errFile.Close()
	e1 := fl.warnFile.Close()
	if e1 != nil && err == nil {
		err = e1
	}
	e1 = fl.infoFile.Close()
	if e1 != nil && err == nil {
		err = e1
	}
	e1 = fl.debugFile.Close()
	if e1 != nil && err == nil {
		err = e1
	}
	return err
}

func newFileLogger(rootPath string) fileLogger {
	var fl fileLogger
	switch currentToFile {
	case TF_ONE:
		fl = &oneFileLogger{
			rootPath: rootPath,
		}
	case TF_SEPARATE:
		fl = &separateFileLogger{
			rootPath: rootPath,
		}
	default:
		panic(fmt.Sprintf("no support for %d", currentToFile))
	}
	fl.NewFile()
	return fl
}

func getFileLogger() fileLogger {
	fileLoggerInstanceOnce.Do(func() {
		root, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		rootPath := filepath.Join(root, "log")
		err = os.MkdirAll(rootPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		fileLoggerInstance = newFileLogger(rootPath)
	})
	if fileLoggerInstance.FileTime().YearDay() != time.Now().YearDay() {
		fileLoggerInstance.NewFile()
	}
	return fileLoggerInstance
}

func writeToFile(message []byte, level loggerMessageLevel) {
	if currentToFile == TF_NO {
		return
	}
	ins := getFileLogger()
	ins.Log(message, level)
}

type LogFile interface {
	Write(message []byte)
	NewFile(path string) error
	Close() error
}

type logFile struct {
	file *os.File
	path string
}

func (lf *logFile) NewFile(path string) error {
	f, err := os.OpenFile(lf.path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		writeRuntimeMessage()
		os.Stdout.Write([]byte(wrapMessage(Format(err), LML_ERR)))
		return err
	}
	lf.file = f
	return nil
}

func (lf *logFile) Write(message []byte) {
	if lf.file == nil {
		err := lf.NewFile(lf.path)
		if err != nil {
			return
		}
	}
	lf.file.Write(message)
}

func (lf *logFile) Close() error {
	if lf.file == nil {
		return nil
	}
	err := lf.file.Close()
	if err != nil {
		writeRuntimeMessage()
		os.Stdout.Write([]byte(wrapMessage(Format(err), LML_ERR)))
		return err
	}
	lf.file = nil
	return nil
}

func SetToFile(to toFile) {
	if fileLoggerInstance == nil {
		currentToFile = to
		return
	}
	err := fileLoggerInstance.Close()
	if err != nil {
		writeRuntimeMessage()
		os.Stdout.Write([]byte(wrapMessage(Format(err), LML_ERR)))
	}
	root, _ := os.Getwd()
	rootPath := filepath.Join(root, "log")
	fileLoggerInstance = newFileLogger(rootPath)
}
