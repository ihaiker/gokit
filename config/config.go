package config

import (
	fileKit "github.com/ihaiker/gokit/files"
	"github.com/jinzhu/configor"
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

	extends := []string{"toml", "yaml", "json"}
	for _, ext := range extends {
		paths := GetStandardConfigurationLocation(name, module, ext)
		for _, path := range paths {
			if fileKit.IsExistFile(path) {
				reg.AddPath(path)
			}
		}
	}
	return reg
}

func (this *register) SearchConfigLocations() []string {
	paths := []string{}
	paths = append(paths, GetStandardConfigurationLocation(this.name, this.module, "json")...)
	paths = append(paths, GetStandardConfigurationLocation(this.name, this.module, "yaml")...)
	paths = append(paths, GetStandardConfigurationLocation(this.name, this.module, "toml")...)
	return paths
}

func (this *register) With(config *configor.Config) *register {
	this.config = config
	return this;
}

func (this *register) AddPath(path ...string) *register {
	this.path = append(this.path, path...)
	return this
}

func (this *register) Marshal(cfg interface{}) error {
	return configor.New(this.config).Load(cfg, this.path...)
}
