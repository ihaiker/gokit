package ssh

import (
	"fmt"
	"net"
	"log"
	"io"
	"sync"
	"golang.org/x/crypto/ssh"
)
/*
	ssh -R/-D/-L  remote:port:local:port
	-R 本地端口开到服务器上
	-L 服务器端口开到本地
*/

type SSHTunnel struct {
	remote, local string
	ssh           *ssh.Client
	quit          chan bool

	listener      net.Listener
	closeOne      sync.Once
}

func (self *SSHTunnel) Name() string {
	return fmt.Sprintf("[%s-%s]", self.local, self.remote)
}

func (self *SSHTunnel) Close() {
	log.Println("close tunnel ", self.Name())

	if self.listener != nil {
		self.listener.Close()
	}

}

//通道是否是激活状态
func (self *SSHTunnel) IsActive() bool {
	return self.listener != nil
}

//激活通道
func (self *SSHTunnel) Activate() error {
	if self.IsActive() {
		return nil
	}
	if err := self.Test(); err != nil {
		return err
	}

	log.Println("create net listener, local endpoint that will forward.", self.local)
	localListenr, err := net.Listen("tcp", self.local)
	if err != nil {
		return err
	}else{
		self.listener = localListenr
	}

	//开始监听端口
	go func() {
		for {
			local, err := localListenr.Accept()
			if err != nil {
				log.Println("localListenr error: ", err)
				return
			}
			log.Println("accept :", local.RemoteAddr().String())
			go func() {
				//拨通远程
				remote, err := self.ssh.Dial("tcp", self.remote)
				if err != nil {
					log.Println("dial remote", self.remote, "error :", err)
					local.Close()
					local = nil
					return
				}
				quit := make(chan bool)
				go iocopy(local, remote, quit)
				go iocopy(remote, local, quit)
				for {
					select {
					case <-quit:
						log.Println("close the tunnel.")
						if remote != nil {
							remote.Close()
							remote = nil
						}
						if local != nil {
							local.Close()
							local = nil
						}
					}
				}
			}()
		}
	}()

	return nil
}

//测试远程和本地连接是否可以使用，
func (self *SSHTunnel) Test() error {
	if localListenr, err := net.Listen("tcp", self.local); err != nil {
		log.Println("test tcp listen", self.local, "error.", err.Error())
		return err
	} else {
		log.Println("test tcp listen", self.local, "ok.")
		localListenr.Close()
	}

	if remoteTCP, err := self.ssh.Dial("tcp", self.remote); err != nil {
		log.Println("test dial tcp", self.remote, "error. ", err.Error())
		return err
	} else {
		log.Println("test dial tcp", self.remote, "ok. ")
		remoteTCP.Close()
	}
	return nil
}

func iocopy(writer, reader net.Conn, quit chan bool) {
	_, err := io.Copy(writer, reader)
	if err != nil && err != io.EOF {
		log.Printf("io.Copy error: %s", err)
	}
	quit <- true
}

func NewTunnel(ssh *ssh.Client, remote, local string) *SSHTunnel {
	tunnel := &SSHTunnel{
		remote:remote,
		local:local,
		ssh:ssh,
		quit:make(chan bool, 2),
	}
	return tunnel
}