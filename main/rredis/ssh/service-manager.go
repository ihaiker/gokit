package ssh

import (
	"log"
	"github.com/ihaiker/gokit/main/rredis/config"
)


//ssh服务管理器
type SSHManager struct {
	sshServices map[string]*SSHService
}

func NewSSHManager() *SSHManager {
	return &SSHManager{sshServices:make(map[string]*SSHService)}
}

//获取service并且自动连接
func (self *SSHManager) Get(redis *config.Redis) (*SSHService, error) {
	name := redis.SSHConfig.User + "@" + redis.SSHConfig.Host
	if sshService, has := self.sshServices[name]; has {
		log.Printf("ssh service %s found", name)
		if err := sshService.Connect(); err != nil {
			return nil, err
		}
		return sshService, nil
	} else {
		log.Printf("ssh service %s not found, to create one.", name)
		sshService = NewSSHService(redis.SSHConfig)
		if err := sshService.Connect(); err != nil {
			return nil, err
		} else {
			self.sshServices[name] = sshService
			return sshService, err
		}
	}
}

//关闭所有连接
func (self *SSHManager) CloseAll() {
	log.Println("close all ssh serivce")
	for name, ss := range self.sshServices {
		ss.Close()
		delete(self.sshServices, name)
	}
}

