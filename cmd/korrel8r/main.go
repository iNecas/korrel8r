// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/korrel8r/korrel8r/internal/pkg/build"
	"github.com/korrel8r/korrel8r/internal/pkg/logging"
	"github.com/korrel8r/korrel8r/internal/pkg/must"
	"github.com/korrel8r/korrel8r/internal/pkg/test/mock"
	"github.com/korrel8r/korrel8r/pkg/config"
	"github.com/korrel8r/korrel8r/pkg/domains/alert"
	"github.com/korrel8r/korrel8r/pkg/domains/incident"
	"github.com/korrel8r/korrel8r/pkg/domains/k8s"
	logdomain "github.com/korrel8r/korrel8r/pkg/domains/log"
	"github.com/korrel8r/korrel8r/pkg/domains/metric"
	"github.com/korrel8r/korrel8r/pkg/domains/netflow"
	"github.com/korrel8r/korrel8r/pkg/domains/trace"
	"github.com/korrel8r/korrel8r/pkg/engine"
	"github.com/korrel8r/korrel8r/pkg/engine/traverse"
	"github.com/spf13/cobra"
)

var (
	log = logging.Log()

	rootCmd = &cobra.Command{
		Use:     "korrel8r",
		Short:   "REST service to correlate observability data",
		Version: build.Version,
	}

	// Global Flags
	outputFlag  = rootCmd.PersistentFlags().StringP("output", "o", "yaml", "Output format: [json, json-pretty, yaml]")
	verboseFlag = rootCmd.PersistentFlags().IntP("verbose", "v", 0, "Verbosity for logging (0: notice/error/warn, 1: info, 2: debug, 3: trace-per-request, 4: trace-per-rule, 5: trace-per-object)")
	configFlag  = rootCmd.PersistentFlags().StringP("config", "c", getConfig(), "Configuration file")
	panicFlag   = rootCmd.PersistentFlags().Bool("panic", false, "Panic on error")
	// TODO: remove sync search once async search is full tested.
	syncFlag = rootCmd.PersistentFlags().Bool("sync", false, "Deprectated: synchronous search, will be removed")
	// see profile.go for profile flag
)

const (
	configEnv     = "KORREL8R_CONFIG"
	defaultConfig = "/etc/korrel8r/korrel8r.yaml"
)

var profileStop interface{ Stop() }

func init() {
	_ = rootCmd.PersistentFlags().MarkHidden("panic")
	_ = rootCmd.PersistentFlags().MarkHidden("sync")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	cobra.OnInitialize(func() {
		logging.Init(verboseFlag)
		k8s.SetLogger(logging.Log())
		if profileFlag != nil {
			profileStop = StartProfile()
		}
	})

	cobra.OnFinalize(func() {
		if profileStop != nil {
			profileStop.Stop()
		}
	})
}

// getConfig looks for the default configuration file.
func getConfig() string {
	if config, ok := os.LookupEnv(configEnv); ok {
		return config // Use env. var. if set.
	}
	return defaultConfig
}

func main() {
	defer func() {
		// Code in this package will panic with an error to cause an exit.
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || *panicFlag {
				panic(r)
			}
			fmt.Fprintf(os.Stderr, "\n%v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	must.Must(rootCmd.Execute())
}

func newEngine() (*engine.Engine, config.Configs) {
	traverse.New = traverse.NewAsync // Default to async
	if *syncFlag {
		traverse.New = traverse.NewSync
	}
	c := must.Must1(config.Load(*configFlag))
	e := must.Must1(engine.Build().
		Domains(k8s.Domain, logdomain.Domain, netflow.Domain, trace.Domain, alert.Domain, metric.Domain, incident.Domain, mock.Domain("mock")).
		Config(c).
		Engine())
	return e, c
}
