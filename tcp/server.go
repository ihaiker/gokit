package gotcp

import (
	"net"
	"sync"
	"time"
)


type Server struct {
	config         *Config         // server configuration
	callback       ConnCallback    // message callbacks in connection
	protocol       Protocol        // customize packet protocol
	exitChan       chan struct{}   // notify all goRoutines to shutdown
	waitGroup      *sync.WaitGroup // wait for all goRoutines
	heartbeatMaker HeartbeatHandlerMaker
}

// NewServer creates a server
func NewServer(callback ConnCallback, protocol Protocol) *Server {
	return &Server{
		config:    DefConfig(),
		callback:  callback,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

func (self *Server) SetConfig(cfg *Config) {
	self.config = cfg
}

func (self *Server) SetHeartbeatMaker(maker HeartbeatHandlerMaker) {
	self.heartbeatMaker = maker
}

// Start starts service
func (s *Server) Start(listener *net.TCPListener, acceptTimeout time.Duration) {
	s.waitGroup.Add(1)
	defer func() {
		listener.Close()
		s.waitGroup.Done()
	}()

	for {
		select {
		case <-s.exitChan:
			return
		default:
		}

		listener.SetDeadline(time.Now().Add(acceptTimeout))

		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		s.waitGroup.Add(1)
		go func() {
			var hb *Heartbeat
			if s.heartbeatMaker != nil {
				hb = s.heartbeatMaker(conn)
			}
			newConn(conn, s, hb).Do()
			s.waitGroup.Done()
		}()
	}
}

// Stop stops service
func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}
