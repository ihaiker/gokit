package rredis

import (
	"strings"
	"os"
	"net"
	"fmt"
	"strconv"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/ihaiker/gokit/main/rredis/ssh"
	"github.com/ihaiker/gokit/main/rredis/config"
	"errors"
	"log"
)

type Service struct {
	name   string //当前服务REDIS名称
	redis  *redis.Client
	ssh    *ssh.SSHService
	tunnel *ssh.SSHTunnel
	config *config.Redis
}

//关闭连接
func (self *Service) Close() {
	//关闭redis连接
	if self.redisIsActive() {
		self.redis.Close()
		self.redis = nil
	}
	//关闭通道，因为通道只有自己在使用，所以可以关闭
	self.ssh.CloseTunnel(self.tunnel)
}

//是否活跃状态
func (self *Service) isActive() bool {

	return self.redisIsActive() && self.isTunnelActive() && self.sshIsActive()
}

//tunnel is active
func (self *Service) isTunnelActive() bool {
	if self.tunnel == nil {
		return false
	}
	return self.tunnel.IsActive()
}

//redis client 是否可用
func (self *Service) redisIsActive() bool {
	if self.redis != nil {
		resp := self.redis.Cmd("PING")
		if resp.Err == nil {
			return true
		}
		self.redis = nil
	}
	return false
}

func (self *Service) sshIsActive() bool {
	if self.ssh != nil {
		return self.ssh.IsConnected()
	}
	return false
}

//激活通道
func (self *Service) Activate() error {
	if !self.isActive() {
		//ssh激活
		if err := self.ssh.Connect(); err != nil {
			return err
		}
		local := fmt.Sprintf("127.0.0.1:%d", randomPort())
		//激活tunnel
		if self.tunnel == nil {
			remote := self.config.Host + ":" + strconv.Itoa(self.config.Post)
			if tunnel, err := self.ssh.CreateTunnel(remote, local); err != nil {
				return err
			} else {
				self.tunnel = tunnel
			}
		}
		if err := self.tunnel.Activate(); err != nil {
			return err
		}
		//激活redisclient
		if client, err := redis.Dial("tcp", local); err != nil {
			return err
		} else {
			//auth redis
			if self.config.Password != "" {
				resp := client.Cmd("AUTH", self.config.Password)
				if err := resp.Err; err != nil {
					client.Close()
					return errors.New("auth " + self.config.Host + " error." + err.Error())
				}
			}
			self.redis = client
		}
	}
	return nil
}
//执行
func (self *Service) Execute(cmd string) (err error) {
	defer func() {
		if errInfo := recover(); errInfo != nil {
			log.Println(errInfo)
			err = errors.New("execute to error.")
		}
	}()

	if err = self.Activate(); err != nil {
		return
	}

	args := strings.Split(cmd, " ")
	resp := self.redis.Cmd(args[0], args[1:])
	if err = resp.Err; err != nil {
		return
	} else {
		resp.WriteTo(os.Stdout)
		return
	}
}

func randomPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port
	return port
}

func NewService(redis *config.Redis) *Service {
	return &Service{
		name:redis.Name,
		ssh:ssh.NewSSHService(redis.SSHConfig),
		config:redis,
	}
}
