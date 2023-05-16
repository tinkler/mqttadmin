/*
build the go file to binary file and run it
*/
package rerun

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/tinkler/mqttadmin/pkg/parser"
)

type Runner struct {
	runCtx  context.Context
	stopRun context.CancelFunc
	running int32
	waiting sync.Mutex
	enable  bool
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Init() {
	if !r.enable {
		return
	}
	r.runCtx, r.stopRun = context.WithCancel(context.Background())
}

func (r *Runner) Run() {
	if !r.enable {
		return
	}
	r.run(r.runCtx)
}

func (r *Runner) Rerun() {
	if !r.enable {
		return
	}
	r.stopRun()
	r.runCtx, r.stopRun = context.WithCancel(context.Background())
	r.run(r.runCtx)
}

func (r *Runner) Stop() {
	if !r.enable {
		return
	}
	r.stopRun()
}

func (r *Runner) Enable() {
	r.enable = true
}

func (r *Runner) run(ctx context.Context) {
	r.waiting.Lock()
	if !atomic.CompareAndSwapInt32(&r.running, 0, 1) {
		r.waiting.Unlock()
		return
	}
	defer func() {
		atomic.StoreInt32(&r.running, 0)
		r.waiting.Unlock()
	}()
	// get the current dir
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	modulePath := parser.GetModulePath()
	// get the go file name
	goFileName := filepath.Base(modulePath)
	if os.Getenv("GOOS") == "windows" {
		goFileName += ".exe"
	}
	// build the go file to binary file
	build := exec.CommandContext(ctx, "go", "build")
	build.Dir = dir
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	err = build.Run()
	if err != nil {
		if ctx.Err() != nil {
			return
		} else {
			panic(err)
		}
	}
	if ctx.Err() != nil {
		return
	}
	// run the binary file
	run := exec.CommandContext(ctx, "./"+goFileName)
	run.Dir = dir
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	err = run.Start()
	if err != nil {
		if ctx.Err() != nil {
			return
		} else {
			panic(err)
		}
	}
	fmt.Println("rerun success")
}
