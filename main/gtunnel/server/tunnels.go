package server

import (
	"github.com/ihaiker/gokit/main/gtunnel/config"
	"net"
	"log"
	"io"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
	"github.com/ihaiker/gokit/files"
)

type GTunnel struct {
	serverConfig *config.ServerConfig
	tunnelConfig *config.TunnelConfig

	sshClient    *ssh.Client
	localListenr net.Listener
	closeOne     *sync.Once
}

func StartTunnel(serverConfig *config.ServerConfig, tunnelConfig *config.TunnelConfig) (*GTunnel, error) {
	tunnel := &GTunnel{
		serverConfig:serverConfig,
		tunnelConfig:tunnelConfig,
		closeOne:new(sync.Once),
	}
	return tunnel, tunnel.Start()
}

func (self *GTunnel) connect() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: self.serverConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(self.serverConfig.Password),
		},
	}
	if self.serverConfig.Privatekey != "" {
		privateKey, err := fileKit.New(self.serverConfig.Privatekey).ToBytes()
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		config = &ssh.ClientConfig{
			User: self.serverConfig.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
		}
	}

	if self.serverConfig.DialTimeoutSecond > 0 {
		connNet, err := net.DialTimeout("tcp", self.serverConfig.Address, time.Duration(self.serverConfig.DialTimeoutSecond) * time.Second)
		if err != nil {
			return nil, err
		}
		sc, chans, reqs, err := ssh.NewClientConn(connNet, self.serverConfig.Address, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(sc, chans, reqs), nil
	} else {
		return ssh.Dial("tcp", self.serverConfig.Address, config)
	}
}

func (self *GTunnel) Start() error {
	//本地监听服务,并且测试是否可以绑定
	ln, err := net.Listen("tcp", self.tunnelConfig.Local)
	if err != nil {
		return err
	}
	self.localListenr = ln

	//链接SSH服务
	sshClient, err := self.connect()
	if err != nil {
		self.Stop()
		return err
	}
	self.sshClient = sshClient

	//测试链接服务的接口如果不能链接就不能创建成功
	if sshConn, err := sshClient.Dial("tcp", self.tunnelConfig.Address); err != nil {
		self.Stop()
		return err
	} else {
		sshConn.Close()//把测试链接关闭
	}
	go self.accept()
	return nil
}

func (self *GTunnel) accept() {
	for {
		conn, err := self.localListenr.Accept()
		if err != nil {
			if !self.tunnelConfig.Active {
				return
			}
			log.Println(err)
		}
		if err := self.handleConnectionAndForward(conn); err != nil {
			log.Println("handle connect error ", err)
		}
	}
}

func (self *GTunnel) handleConnectionAndForward(conn net.Conn) error {
	sshConn, err := self.sshClient.Dial("tcp", self.tunnelConfig.Address)
	if err != nil {
		return err
	}
	go func() {
		if _, err := io.Copy(conn, sshConn); err != nil && self.tunnelConfig.Active {
			log.Println(err)
		}
		self.Stop()
	}()
	go func() {
		if _, err := io.Copy(sshConn, conn); err != nil && self.tunnelConfig.Active {
			log.Println(err)
		}
		self.Stop();
	}()
	return nil
}

func (self *GTunnel) Stop() error {
	self.closeOne.Do(func() {
		log.Println("stop the tunnel" + self.tunnelConfig.Name)
		self.tunnelConfig.Active = false
		if self.localListenr != nil {
			log.Println("[info]", "close listener")
			self.localListenr.Close()
		}
		if self.sshClient != nil {
			log.Println("[info]", "ssh client")
			self.sshClient.Close()
		}
	})
	return nil
}