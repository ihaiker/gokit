package config

import (
	"os"
	"github.com/ihaiker/gokit/files"
	"path/filepath"
	"fmt"
	"github.com/ihaiker/gokit/config/json"
)

//redis配置
type Redis struct {
	Name     string `json:"name"`
	Host     string	`json:"host"`
	Post     int	`json:"port"`
	Password string `json:"password"`
	SSHConfig *SSHConfig 	`json:"ssh"`
}
//ssh配置
type SSHConfig struct {
	Host              string
	User              string
	Password          string
	Privatekey        string
	DialTimeoutSecond int
	MaxDataThroughput uint64
}
type Config struct {
	HistoryPath string `json:"history"`
	SSHConfig   *SSHConfig `json:"ssh"`
	RemoteRedis []*Redis `json:"remote"`
}

func ReadConfig(args []string) *Config {
	cfgFile := ""
	if len(args) == 2 {
		cfgFile = args[1]
	}
	if cfgFile == "" {
		cfgFile1 := os.Getenv("HOME") + fileKit.Separator + ".rredis.json"
		cfgFile2 := filepath.Dir(args[0]) + fileKit.Separator + ".rredis.json"
		if !fileKit.IsExistFile(cfgFile1) && !fileKit.IsExistFile(cfgFile2) {
			fmt.Println("the config file not found ! \nat " + cfgFile1 + "\nat " + cfgFile2)
			os.Exit(1)
		}
		if fileKit.IsExistFile(cfgFile1) {
			cfgFile = cfgFile1
		} else if fileKit.IsExistFile(cfgFile2) {
			cfgFile = cfgFile2
		}
	}
	if !fileKit.IsExistFile(cfgFile) {
		fmt.Println("can't found the config " + cfgFile)
		os.Exit(1)
	}
	config := &Config{}
	if jsonCfg, err := json.Config(fileKit.New(cfgFile)); err != nil {
		panic(err)
	} else if err := jsonCfg.Unmarshal(config); err != nil {
		panic(err)
	}

	if config.HistoryPath == "" {
		config.HistoryPath = os.Getenv("HOME") + fileKit.Separator + ".rredis_history"
	}

	//设置每个Redis使用的服务
	for _, r := range config.RemoteRedis {
		if r.SSHConfig == nil {
			r.SSHConfig = config.SSHConfig
		}
	}

	return config
}
