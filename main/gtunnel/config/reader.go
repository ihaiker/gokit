package config

import (
	"os"
	"github.com/ihaiker/gokit/files"
	"path/filepath"
	"github.com/ihaiker/gokit/config/json"
	"errors"
)

const CONFIG_FILE_NAME = ".gtunnel.json"

func ReadConfig(args []string) (*Config,error) {
	cfgFile := ""
	if len(args) == 1 {
		cfgFile = args[0]
	}
	if cfgFile == "" {
		cfgFile1 := os.Getenv("HOME") + fileKit.Separator + CONFIG_FILE_NAME
		cfgFile2 := filepath.Dir(os.Args[0]) + fileKit.Separator + CONFIG_FILE_NAME
		if !fileKit.IsExistFile(cfgFile1) && !fileKit.IsExistFile(cfgFile2) {
			return nil,errors.New("the config file not found ! \nat " + cfgFile1 + "\nat " + cfgFile2)
		}
		if fileKit.IsExistFile(cfgFile1) {
			cfgFile = cfgFile1
		} else if fileKit.IsExistFile(cfgFile2) {
			cfgFile = cfgFile2
		}
	}
	if !fileKit.IsExistFile(cfgFile) {
		return nil,errors.New("can't found the config " + cfgFile)
	}
	config := &Config{}
	if jsonCfg, err := json.Config(fileKit.New(cfgFile)); err != nil {
		panic(err)
	} else if err := jsonCfg.Unmarshal(config); err != nil {
		panic(err)
	}
	return config,nil
}
