package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/baudtime/agent/plugin/manager"
	"github.com/baudtime/agent/vars"
	"github.com/go-kit/kit/log/level"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	vars.Init()

	if vars.Cfg.ServicePort != nil {
		go func() {
			http.ListenAndServe(fmt.Sprintf(":%d", *vars.Cfg.ServicePort), nil)
		}()
	}

	mgr, err := manager.New(vars.Cfg.Inputs.Enabled, vars.Cfg.Processors.Enabled, vars.Cfg.Aggregators.Enabled, vars.Cfg.Outputs.Enabled)
	if err != nil {
		level.Error(vars.Logger).Log("msg", "failed to init plugins manager", "err", err)
		return
	}

	mgr.Start()
	defer mgr.Stop()

	waitFor(syscall.SIGTERM, syscall.SIGINT)
}

func waitFor(sig ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	for range c {
		return
	}
}
