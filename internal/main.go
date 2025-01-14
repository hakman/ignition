// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/flatcar-linux/ignition/internal/exec"
	"github.com/flatcar-linux/ignition/internal/exec/stages"
	_ "github.com/flatcar-linux/ignition/internal/exec/stages/disks"
	_ "github.com/flatcar-linux/ignition/internal/exec/stages/fetch"
	_ "github.com/flatcar-linux/ignition/internal/exec/stages/files"
	"github.com/flatcar-linux/ignition/internal/log"
	"github.com/flatcar-linux/ignition/internal/oem"
	"github.com/flatcar-linux/ignition/internal/version"
)

func main() {
	flags := struct {
		clearCache   bool
		configCache  string
		fetchTimeout time.Duration
		oem          oem.Name
		root         string
		stage        stages.Name
		version      bool
		logToStdout  bool
	}{}

	flag.BoolVar(&flags.clearCache, "clear-cache", false, "clear any cached config")
	flag.StringVar(&flags.configCache, "config-cache", "/run/ignition.json", "where to cache the config")
	flag.DurationVar(&flags.fetchTimeout, "fetch-timeout", exec.DefaultFetchTimeout, "initial duration for which to wait for config")
	flag.Var(&flags.oem, "oem", fmt.Sprintf("current oem. %v", oem.Names()))
	flag.StringVar(&flags.root, "root", "/", "root of the filesystem")
	flag.Var(&flags.stage, "stage", fmt.Sprintf("execution stage. %v", stages.Names()))
	flag.BoolVar(&flags.version, "version", false, "print the version and exit")
	flag.BoolVar(&flags.logToStdout, "log-to-stdout", false, "log to stdout instead of the system log when set")

	flag.Parse()

	if flags.version {
		fmt.Printf("%s\n", version.String)
		return
	}

	if flags.oem == "" {
		fmt.Fprint(os.Stderr, "'--oem' must be provided\n")
		os.Exit(2)
	}

	if flags.stage == "" {
		fmt.Fprint(os.Stderr, "'--stage' must be provided\n")
		os.Exit(2)
	}

	logger := log.New(flags.logToStdout)
	defer logger.Close()

	logger.Info(version.String)
	logger.Info("Stage: %v", flags.stage)

	if flags.clearCache {
		if err := os.Remove(flags.configCache); err != nil {
			logger.Err("unable to clear cache: %v", err)
		}
	}

	oemConfig := oem.MustGet(flags.oem.String())
	fetcher, err := oemConfig.NewFetcherFunc()(&logger)
	if err != nil {
		logger.Crit("failed to generate fetcher: %s", err)
		os.Exit(3)
	}
	engine := exec.Engine{
		Root:         flags.root,
		FetchTimeout: flags.fetchTimeout,
		Logger:       &logger,
		ConfigCache:  flags.configCache,
		OEMConfig:    oemConfig,
		Fetcher:      &fetcher,
	}

	err = engine.Run(flags.stage.String())
	if statusErr := engine.OEMConfig.Status(flags.stage.String(), *engine.Fetcher, err); statusErr != nil {
		logger.Err("POST Status error: %v", statusErr.Error())
	}
	if err != nil {
		logger.Crit("Ignition failed: %v", err.Error())
		os.Exit(1)
	}
	logger.Info("Ignition finished successfully")
}
