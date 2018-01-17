package tcpKit

import (
    "sync"
    "net"
    "time"
    "github.com/ihaiker/gokit/commons/logs"
    "strings"
)

type Server struct {
    config *Config // server configuration
    maker struct {
        handler  HandlerMaker  // message callbacks in connection
        protocol ProtocolMaker // customize packet protocol
    }
    clients   map[string]*Connect
    exitChan  chan struct{}   // notify all goroutines to shutdown
    waitGroup *sync.WaitGroup // wait for all goroutines
    logger    *logs.LoggerEntry
}

func (s *Server) StartAt(addr string) error {
    if tcpAddr, err := net.ResolveTCPAddr("tcp4", addr); err != nil {
        return err
    } else if listener, err := net.ListenTCP("tcp", tcpAddr); err != nil {
        return err
    } else {
        s.Start(listener)
        return nil
    }
}

// Start starts service
func (s *Server) Start(listener *net.TCPListener) {
    s.waitGroup.Add(1)
    defer func() {
        listener.Close()
        s.waitGroup.Done()
    }()
    s.logger.Debugf("TCP服务启动：%s", listener.Addr().String())
    for {
        select {
        case <-s.exitChan:
            return
        default:
        }
        listener.SetDeadline(time.Now().Add(s.config.AcceptTimeout))
        conn, err := listener.AcceptTCP()
        if err != nil {
            if strings.Contains(err.Error(), "i/o timeout") {
                continue
            }
            s.logger.Errorf("服务监听错误：%s", err)
        }
        s.waitGroup.Add(1)
        go func() {
            addr := conn.RemoteAddr().String()
            s.logger.Debugf("客户端连接服务器：%s", addr)
            c := newConnect(s, conn)
            s.clients[addr] = c
            c.Do(func(c *Connect) {
                s.logger.Debugf("客户端关闭连接：%s", addr)
                delete(s.clients, addr)
                s.waitGroup.Done()
            })
        }()
    }
}

// Stop stops service
func (s *Server) Stop() {
    s.logger.Debug("关闭TCP服务")
    close(s.exitChan)
    for _, v := range s.clients {
        if v != nil {
            v.Close()
        }
    }
    s.waitGroup.Wait()
    s.logger.Debug("关闭TCP服务完成")
}

// NewServer creates a server
func NewServer(config *Config, handlerMaker HandlerMaker, protocolMaker ProtocolMaker) *Server {
    if config.IdleTimeout == 0 {
        config.IdleTimeout = config.IdleTime
    }
    return &Server{
        config: config,
        maker: struct {
            handler  HandlerMaker
            protocol ProtocolMaker
        }{handler: handlerMaker, protocol: protocolMaker},
        clients:   make(map[string]*Connect),
        exitChan:  make(chan struct{}),
        waitGroup: &sync.WaitGroup{},
        logger:    logs.Logger("tcpKit"),
    }
}
