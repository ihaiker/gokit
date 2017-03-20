package ssh

import (
	"golang.org/x/crypto/ssh"
	"github.com/ihaiker/gokit/commons/time"
	"github.com/ihaiker/gokit/files"
	"log"
	"time"
	"github.com/ihaiker/gokit/main/rredis/config"
)

type SSHService struct {
	cfg     *config.SSHConfig
	ssh     *ssh.Client

	tunnels map[string]*SSHTunnel
}

func NewSSHService(cfg *config.SSHConfig) *SSHService {
	return &SSHService{
		cfg        : cfg,
		tunnels : make(map[string]*SSHTunnel)}
}

func (self *SSHService) Test() error {
	if err := self.Connect(); err != nil {
		return err
	} else {
		self.ssh.Close()
		self.ssh = nil
		return nil
	}
}

//检测是否连接状态
func (self *SSHService) IsConnected() bool {
	if self.ssh != nil {
		return true
	}
	return false
}

//链接ssh服务器
func (self *SSHService) Connect() error {
	name := self.cfg.User + "@" + self.cfg.Host
	if self.IsConnected() {
		log.Println("ssh service", name, "is connected.")
		return nil
	}
	log.Println("to connect the ssh service : ", name)

	sshConfig := &ssh.ClientConfig{
		User: self.cfg.User,
	}
	sshConfig.Timeout = timeKit.Duration(time.Millisecond, self.cfg.DialTimeoutSecond)
	if self.cfg.Password != "" {
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.Password(self.cfg.Password),
		}
	} else {
		pkBody, err := fileKit.New(self.cfg.Privatekey).ToBytes()
		if err != nil {
			return err
		}
		pk, err := ssh.ParsePrivateKey(pkBody)
		if err != nil {
			return err
		}
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(pk),
		}
	}
	if client, err := ssh.Dial("tcp", self.cfg.Host, sshConfig); err != nil {
		return err
	} else {
		log.Println("connect ssh service ", name, " successed.")
		self.ssh = client
		return nil
	}
}

//创建服务通道
func (self *SSHService) CreateTunnel(remoteAddr, localAddr string) (*SSHTunnel, error) {
	tunnel := NewTunnel(self.ssh, remoteAddr, localAddr)
	log.Printf("create ssh tunnel %s use server %s\n", tunnel.Name(), self.cfg.Host)
	if tu, has := self.tunnels[tunnel.Name()]; has {
		return tu, nil
	} else {
		if err := tunnel.Activate(); err != nil {
			return nil, err
		}
		self.tunnels[tunnel.Name()] = tunnel
		return tunnel, nil
	}
}

//关闭通道，如果所有通道关闭，服务将关闭
func (self *SSHService) CloseTunnel(tunnel *SSHTunnel) {
	if tunnel != nil {
		tunnel.Close()
	}
}

func (self *SSHService) Close() {
	//关闭此服务下的所有tunnel
	if self.tunnels != nil {
		log.Println("close tunnels at ",self.name())
		for _,t := range  self.tunnels{
			t.Close()
		}
	}
	//关闭服务
	if self.ssh != nil {
		log.Println("to close ssh serivce", self.name())
		self.ssh.Close()
	}
}

func (self *SSHService) name() string{
	name := self.cfg.User + "@" + self.cfg.Host
	return name
}