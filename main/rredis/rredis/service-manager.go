package rredis

import (
	"errors"
	"log"
	"github.com/ihaiker/gokit/main/rredis/config"
)

type ServiceManager struct {
	cfg           *config.Config

	serivces      map[string]*Service
	currentServer string
}

func (self *ServiceManager) Close() {
	log.Println("close services.")
	for name, service := range self.serivces {
		if service.isActive() {
			service.Close()
			log.Println("close ", name)
		}
	}
}

func (self *ServiceManager) Use(name string) error {
	log.Println("use database :", name, ", the current :", self.currentServer, ".")
	if self.currentServer == name && self.Current().isActive() {
		return nil
	}
	if _, has := self.serivces[name]; !has {
		return errors.New("the database not found!")
	}
	self.currentServer = name
	if !self.Current().isActive() {
		log.Println("to activate ",self.currentServer)
		if err := self.Current().Activate(); err != nil {
			return err
		}
	}else{
		log.Println("the database is activate :",self.currentServer)
	}
	return nil
}

func (self *ServiceManager) Current() *Service {
	s, has := self.serivces[self.currentServer]
	if has {
		return s
	}
	return nil
}

func InitServiceManager(cfg *config.Config) *ServiceManager {
	log.Println("init service manager")
	sm := &ServiceManager{
		cfg:cfg,
		serivces:make(map[string]*Service),
	}
	for _, redis := range cfg.RemoteRedis {
		if redis.SSHConfig == nil {
			redis.SSHConfig = cfg.SSHConfig
		}
		sm.serivces[redis.Name] = NewService(redis)
	}
	return sm
}