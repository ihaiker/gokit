package server

import (
    "github.com/ihaiker/gokit/main/gtunnel/config"
    "net"
    "github.com/ihaiker/gokit/commons/logs"
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
    localListener net.Listener
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
    self.localListener = ln

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
        conn, err := self.localListener.Accept()
        if err != nil {
            if !self.tunnelConfig.Active {
                return
            }
            logs.Error(err)
        }
        if err := self.handleConnectionAndForward(conn); err != nil {
            logs.Error("handle connect error ", err)
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
            logs.Error(err)
        }
        self.Stop()
    }()
    go func() {
        if _, err := io.Copy(sshConn, conn); err != nil && self.tunnelConfig.Active {
            logs.Error(err)
        }
        self.Stop();
    }()
    return nil
}

func (self *GTunnel) Stop() error {
    self.closeOne.Do(func() {
        logs.Infof("stop the tunnel %s", self.tunnelConfig.Name)
        self.tunnelConfig.Active = false
        if self.localListener != nil {
            logs.Debug("close listener ", self.tunnelConfig.Name)
            self.localListener.Close()
        }
        if self.sshClient != nil {
            logs.Debug("ssh client", self.tunnelConfig.Name)
            self.sshClient.Close()
        }
    })
    return nil
}