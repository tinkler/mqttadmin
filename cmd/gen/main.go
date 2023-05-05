/*
live code generator
*/
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/tinkler/mqttadmin/internal/gen"
	"github.com/tinkler/mqttadmin/pkg/parser"
)

// mkdirForOutput
func mkdirForOutput(dir string, codes []*gen.GenCodeConf) {
	for _, c := range codes {
		if c.OutDir == "" {
			continue
		}
		err := os.MkdirAll(c.OutDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

// generate all files when start
func beforeWatch(dir string, codes []*gen.GenCodeConf) (cache map[string]*parser.Package, err error) {
	cache = make(map[string]*parser.Package)
	// traverse dir
	err = filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !info.IsDir() {
			return nil
		}
		if dir == path {
			return nil
		}
		pkg, err := parser.ParsePackage(path)
		if err != nil {
			log.Fatal(err)
		}
		cache[pkg.Name] = pkg
		return nil
	})
	if err != nil {
		return
	}
	for _, pkg := range cache {
		for _, c := range codes {
			switch c.Typ {
			case "ts":
				err = parser.GenerateTSCode(c.OutDir, pkg, cache)
				if err != nil {
					log.Fatal(err)
				}
			case "dart":
				err = parser.GenerateDartCode(c.OutDir, pkg, cache)
				if err != nil {
					log.Fatal(err)
				}
			case "swift":
				err = parser.GenerateSwiftCode(c.OutDir, pkg)
				if err != nil {
					log.Fatal(err)
				}
			case "chi_route":
				err = parser.GenerateChiRoutes(c.OutDir, pkg, cache)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return
}

func main() {
	gen.ParseGenConf()
	gf := gen.GetGenConf()
	mkdirForOutput(gf.Dir, gf.Codes)
	pkgs, err := beforeWatch(gf.Dir, gf.Codes)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
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
					for _, c := range gf.Codes {
						switch c.Typ {
						case "ts":
							err = parser.GenerateTSCode(c.OutDir, pkg, pkgs)
							if err != nil {
								log.Fatal(err)
							}
						case "dart":
							err = parser.GenerateDartCode(c.OutDir, pkg, pkgs)
							if err != nil {
								log.Fatal(err)
							}
						case "swift":
							err = parser.GenerateSwiftCode(c.OutDir, pkg)
							if err != nil {
								log.Fatal(err)
							}
						case "chi_route":
							err = parser.GenerateChiRoutes(c.OutDir, pkg, pkgs)
							if err != nil {
								log.Fatal(err)
							}
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
	err = filepath.Walk(gf.Dir, func(path string, info os.FileInfo, e error) error {
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
