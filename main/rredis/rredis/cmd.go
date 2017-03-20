package rredis

import (
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/main/rredis/config"
)

var QUIT_ERROR = errors.New("exit")

var defaultCommand = func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
	if service := servcieManager.Current(); service != nil {
		if !service.isActive() {
			//auto reconnect
			if err := Commands["use"](servcieManager, cfg, servcieManager.currentServer); err != nil {
				return err
			}
		}
		return service.Execute(cmd)
	}
	return errors.New("Have no choice to the server")
}

var Commands = map[string]func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error{
	"use" : func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
		return servcieManager.Use(cmd)
	},
	"dbs" : func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
		for _, redis := range cfg.RemoteRedis {
			fmt.Println(redis.Name, " (", redis.Host, ")")
		}
		return nil
	},
	"help": func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
		return nil
	},
	"exit": func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
		return QUIT_ERROR
	},
	"quit": func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
		return QUIT_ERROR
	},
	"clear": func(servcieManager *ServiceManager, cfg *config.Config, cmd string) error {
		return errors.New("not support")
	},
}
