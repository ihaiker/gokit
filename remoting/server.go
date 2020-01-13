package remoting

import (
	"github.com/ihaiker/gokit/concurrent/executors"
	"net"
	"sync"
)

type Server interface {
	//启动服务
	Start() error

	//关闭服务
	Stop() error

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

	worker executors.ExecutorService
}

func NewServer(address string, config *Config, handlerMaker HandlerMaker, coderMaker CoderMaker) Server {
	return &tcpServer{
		address: address, config: config,

		handlerMaker: handlerMaker,
		coderMaker:   coderMaker,

		clients: NewIpClientManager(),
		worker:  executors.Fixed(config.AsyncHandlerGroup),

		exitChan: make(chan struct{}),
		closeOne: new(sync.Once), waitGroup: new(sync.WaitGroup),
	}
}

func (s *tcpServer) startAccept() {
	defer s.waitGroup.Done()

	logger.Info("remoting start：", s.listener.Addr().String())

	for {
		select {
		case <-s.exitChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				return
			}

			channel := newChannel(s.config, s.worker, conn)
			logger.Debug("client connect：", channel)

			channel.coder = s.coderMaker(channel)
			channel.handler = s.handlerMaker(channel)

			s.clients.Add(channel)

			s.waitGroup.Add(1)
			go channel.do(func(Channel) {}, func(c Channel) {
				defer s.waitGroup.Done()
				logger.Debug("client close：", c)
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
	logger.Debug("server stoped")
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
func (s *tcpServer) Stop() error {
	s.closeOne.Do(func() {
		logger.Info("close server")
		_ = s.listener.Close()
		s.worker.Shutdown()
		close(s.exitChan)
	})
	return nil
}
