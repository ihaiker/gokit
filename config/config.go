package config

import (
	"errors"
	"github.com/ihaiker/gokit/files"
	"github.com/jinzhu/configor"
)

var (
	ErrConfigNotFound = errors.New("the config not found")
)

type register struct {
	name   string
	module string
	config *configor.Config
	path   []string
}

func NewConfigRegister(name, module string) *register {
	if module == "" {
		module = name
	}
	reg := &register{
		name: name, module: module,
		path: []string{},
	}

	paths := reg.searchConfigLocations()
	for _, path := range paths {
		if files.IsExistFile(path) {
			reg.AddPath(path)
		}
	}
	return reg
}

func (this *register) searchConfigLocations() []string {
	paths := make([]string, 0)
	paths = append(paths, GetStandardConfigurationLocation(this.name, this.module, "toml")...)
	paths = append(paths, GetStandardConfigurationLocation(this.name, this.module, "yaml")...)
	paths = append(paths, GetStandardConfigurationLocation(this.name, this.module, "json")...)
	return paths
}

func (this *register) With(config *configor.Config) *register {
	this.config = config
	return this
}

func (this *register) AddPath(path ...string) *register {
	this.path = append(path, this.path...)
	return this
}

func (this *register) AddAtPath(at int, path string) *register {
	if len(this.path) > at {
		start := this.path[0:at]
		end := this.path[at:0]
		this.path = append(start, path)
		this.path = append(this.path, end...)
	}
	return this
}

func (this *register) Marshal(cfg interface{}) error {
	if len(this.path) == 0 {
		return ErrConfigNotFound
	}
	return configor.New(this.config).Load(cfg, this.path...)
}
