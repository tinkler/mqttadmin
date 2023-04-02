/*
live code generator
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/tinkler/mqttadmin/pkg/parser"
)

var (
	// listen dir
	dir string
	// listen type
	typ string
)

func init() {
	flag.StringVar(&dir, "dir", "../../pkg/model", "listen dir")
	flag.StringVar(&typ, "type", "ts", "listen type")
	flag.Parse()
}

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if filepath.Ext(event.Name) != ".go" {
						continue
					}
					if strings.HasSuffix(event.Name, "_test.go") {
						continue
					}
					fmt.Println("modified file:", event.Name)

					pkg, err := parser.ParsePackage(filepath.Dir(event.Name))
					if err != nil {
						log.Fatal(err)
					}
					switch typ {
					case "ts":
						// TODO: cache other package
						err = parser.GenerateTSCode("../../", pkg, nil)
						if err != nil {
							log.Fatal(err)
						}
					case "dart":
						err = parser.GenerateDartCode("../../", pkg)
						if err != nil {
							log.Fatal(err)
						}
					case "swift":
						err = parser.GenerateSwiftCode("../../", pkg)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-done:
				return
			}
		}
	}()
	// traverse dir
	err = filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !info.IsDir() {
			return nil
		}
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	<-quit
	fmt.Println("退出")
	done <- true
}
