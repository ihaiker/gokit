package remoting

import (
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Server interface {
	//启动服务
	Start() error

	//关闭服务
	Stop() Server

	Wait()

	//根据客户端clientId获取客户连接
	GetClientManager() ChannelManager

	SetClientManager(manager ChannelManager)
}

type tcpServer struct {
	config *Config // server configuration

	address  string
	listener net.Listener

	handlerMaker HandlerMaker // message callbacks in connection
	coderMaker   CoderMaker   // customize packet coderMaker

	clients ChannelManager

	exitChan chan struct{}

	closeOne  *sync.Once
	waitGroup *sync.WaitGroup
}

func NewServer(address string, config *Config, handlerMaker HandlerMaker, coderMaker CoderMaker) Server {
	return &tcpServer{
		address: address,
		config:  config,

		handlerMaker: handlerMaker,
		coderMaker:   coderMaker,

		clients: NewIpClientManager(),

		exitChan: make(chan struct{}),
		closeOne: new(sync.Once), waitGroup: new(sync.WaitGroup),
	}
}

func (s *tcpServer) startAccept() {
	defer func() {
		_ = s.listener.Close()
		s.waitGroup.Done()
	}()
	logger.Info("服务启动：", s.listener.Addr().String())

	isTcp := reflect.TypeOf(s.listener).String() == reflect.TypeOf(new(net.TCPListener)).String()
	if ! isTcp {
		s.listener.(*net.UnixListener).SetUnlinkOnClose(true)
	}
	for {
		select {
		case <-s.exitChan:
			return
		default:
			if isTcp {
				_ = s.listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
			} else {
				_ = s.listener.(*net.UnixListener).SetDeadline(time.Now().Add(time.Second))
			}

			conn, err := s.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "i/o timeout") {
					logger.Errorf("服务监听错误：%s", err)
				}
				continue
			}

			s.waitGroup.Add(1)
			addr := conn.RemoteAddr().String()
			logger.Debug("客户端连接服务器：", addr)

			channel := newChannel(s.config, conn)

			channel.coder = s.coderMaker(channel)
			channel.handler = s.handlerMaker(channel)

			s.clients.Add(channel)
			go channel.do(func(Channel) {}, func(c Channel) {
				defer s.waitGroup.Done()
				logger.Debug("客户端关闭连接：", addr)
				s.clients.Remove(c)
			})
		}
	}
}

// Start starts service
func (s *tcpServer) Start() (err error) {
	if s.listener, err = Listen(s.address); err != nil {
		return
	}
	s.waitGroup.Add(1)
	go s.startAccept()
	return nil
}

func (s *tcpServer) Wait() {
	s.waitGroup.Wait()
}

//根据客户端clientId获取客户连接
func (s *tcpServer) GetClientManager() ChannelManager {
	return s.clients
}

func (s *tcpServer) SetClientManager(manager ChannelManager) {
	s.clients.Foreach(func(channel Channel) {
		manager.Add(channel)
	})
	s.clients = manager
}

// Stop stops service
func (s *tcpServer) Stop() Server {
	s.closeOne.Do(func() {
		logger.Info("关闭TCP服务")
		close(s.exitChan)
		s.clients.Foreach(func(channel Channel) {
			channel.Close()
		})
		logger.Debug("关闭TCP服务完成")
	})
	s.waitGroup.Wait()
	return s
}
