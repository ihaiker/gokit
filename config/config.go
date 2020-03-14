package config

import (
	"errors"
	fileKit "github.com/ihaiker/gokit/files"
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
		reg.AddConfig(path)
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

func (this *register) AddConfig(path ...string) *register {
	//覆盖模式从后向前
	for _, p := range path {
		this.path = append([]string{p}, this.path...)
	}
	return this
}

func (this *register) MustExitConfig(files ...string) error {
	for _, f := range files {
		if !fileKit.IsExistFile(f) {
			return ErrConfigNotFound
		}
	}
	this.AddConfig(files...)
	return nil
}

func (this *register) Marshal(cfg interface{}) error {
	return configor.New(this.config).Load(cfg, this.path...)
}
