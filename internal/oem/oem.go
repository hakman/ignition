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

package oem

import (
	"fmt"

	"github.com/flatcar-linux/ignition/internal/log"
	"github.com/flatcar-linux/ignition/internal/providers"
	"github.com/flatcar-linux/ignition/internal/providers/aliyun"
	"github.com/flatcar-linux/ignition/internal/providers/azure"
	"github.com/flatcar-linux/ignition/internal/providers/cloudstack"
	"github.com/flatcar-linux/ignition/internal/providers/digitalocean"
	"github.com/flatcar-linux/ignition/internal/providers/ec2"
	"github.com/flatcar-linux/ignition/internal/providers/file"
	"github.com/flatcar-linux/ignition/internal/providers/gce"
	"github.com/flatcar-linux/ignition/internal/providers/noop"
	"github.com/flatcar-linux/ignition/internal/providers/openstack"
	"github.com/flatcar-linux/ignition/internal/providers/packet"
	"github.com/flatcar-linux/ignition/internal/providers/qemu"
	"github.com/flatcar-linux/ignition/internal/providers/virtualbox"
	"github.com/flatcar-linux/ignition/internal/providers/vmware"
	"github.com/flatcar-linux/ignition/internal/providers/vultr"
	"github.com/flatcar-linux/ignition/internal/registry"
	"github.com/flatcar-linux/ignition/internal/resource"
)

// Config represents a set of options that map to a particular OEM.
type Config struct {
	name       string
	fetch      providers.FuncFetchConfig
	newFetcher providers.FuncNewFetcher
	status     providers.FuncPostStatus
}

func (c Config) Name() string {
	return c.name
}

func (c Config) FetchFunc() providers.FuncFetchConfig {
	return c.fetch
}

func (c Config) NewFetcherFunc() providers.FuncNewFetcher {
	if c.newFetcher != nil {
		return c.newFetcher
	}
	return func(l *log.Logger) (resource.Fetcher, error) {
		return resource.Fetcher{
			Logger: l,
		}, nil
	}
}

// Status takes a Fetcher and the error from Run (from engine)
func (c Config) Status(stageName string, f resource.Fetcher, statusErr error) error {
	if c.status != nil {
		return c.status(stageName, f, statusErr)
	}
	return nil
}

var configs = registry.Create("oem configs")

func init() {
	configs.Register(Config{
		name:  "aliyun",
		fetch: aliyun.FetchConfig,
	})
	configs.Register(Config{
		name:  "azure",
		fetch: azure.FetchConfig,
	})
	configs.Register(Config{
		name:  "cloudsigma",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "cloudstack",
		fetch: cloudstack.FetchConfig,
	})
	configs.Register(Config{
		name:  "digitalocean",
		fetch: digitalocean.FetchConfig,
	})
	configs.Register(Config{
		name:  "brightbox",
		fetch: openstack.FetchConfig,
	})
	configs.Register(Config{
		name:  "openstack",
		fetch: openstack.FetchConfig,
	})
	configs.Register(Config{
		name:       "ec2",
		fetch:      ec2.FetchConfig,
		newFetcher: ec2.NewFetcher,
	})
	configs.Register(Config{
		name:  "exoscale",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "gce",
		fetch: gce.FetchConfig,
	})
	configs.Register(Config{
		name:  "hyperv",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "niftycloud",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:   "packet",
		fetch:  packet.FetchConfig,
		status: packet.PostStatus,
	})
	configs.Register(Config{
		name:  "pxe",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "rackspace",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "rackspace-onmetal",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "vagrant",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "vagrant-virtualbox",
		fetch: virtualbox.FetchConfig,
	})
	configs.Register(Config{
		name:  "virtualbox",
		fetch: virtualbox.FetchConfig,
	})
	configs.Register(Config{
		name:  "vmware",
		fetch: vmware.FetchConfig,
	})
	configs.Register(Config{
		name:  "vultr",
		fetch: vultr.FetchConfig,
	})
	configs.Register(Config{
		name:  "interoute",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "qemu",
		fetch: qemu.FetchConfig,
	})
	configs.Register(Config{
		name:  "file",
		fetch: file.FetchConfig,
	})
	configs.Register(Config{
		name:  "metal",
		fetch: noop.FetchConfig,
	})
}

func Get(name string) (config Config, ok bool) {
	config, ok = configs.Get(name).(Config)
	return
}

func MustGet(name string) Config {
	if config, ok := Get(name); ok {
		return config
	} else {
		panic(fmt.Sprintf("invalid OEM name %q provided", name))
	}
}

func Names() (names []string) {
	return configs.Names()
}
