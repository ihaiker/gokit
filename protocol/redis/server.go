// go-redis-server is a helper library for building server software capable of speaking the redis protocol.
// This could be an alternate implementation of redis, a custom proxy to redis,
// or even a completely different backend capable of "masquerading" its API as a redis database.

package redis

import (
	"fmt"
	"net"
	"github.com/ihaiker/gokit/commons/logs"
)

type Server struct {
	Proto        string
	Address      string // TCP address to listen on, ":6389" if empty
	MonitorChans []chan string
	methods      map[string]HandlerFn
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Address
	if srv.Proto == "" {
		srv.Proto = "tcp"
	}
	if srv.Proto == "unix" && addr == "" {
		addr = "/tmp/redis.sock"
	} else if addr == "" {
		addr = ":6379"
	}
	logs.Infof("the server run at %s %s", srv.Proto, addr)
	
	l, e := net.Listen(srv.Proto, addr)
	if e != nil {
		return e
	}
	return srv.serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) serve(l net.Listener) error {
	defer l.Close()
	srv.MonitorChans = []chan string{}
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		go srv.serveClient(rw)
	}
}

// Serve starts a new redis session, using `conn` as a transport.
// It reads commands using the redis protocol, passes them to `handler`,
// and returns the result.
func (srv *Server) serveClient(conn net.Conn) (err error) {
	var clientAddress string
	
	defer func() {
		if err != nil {
			fmt.Fprintf(conn, "-%s\n", err)
		}
		logs.Debug("the client ",clientAddress, " disconnect.")
		conn.Close()
	}()
	
	switch co := conn.(type) {
	case *net.UnixConn:
		f, err := conn.(*net.UnixConn).File()
		if err != nil {
			return err
		}
		clientAddress = f.Name()
	default:
		clientAddress = co.RemoteAddr().String()
	}
	logs.Debug("the client ", clientAddress," connect.")
	
	for {
		request, err := parseRequest(conn)
		if err != nil {
			return err
		}
		request.Host = clientAddress
		reply, err := srv.Apply(request)
		if err != nil {
			return err
		}
		if _, err = reply.WriteTo(conn); err != nil {
			return err
		}
	}
	return nil
}

func NewServer() *Server {
	srv := &Server{
		Proto:        "tcp",
		Address:         ":6379",
		MonitorChans: []chan string{},
		methods:      make(map[string]HandlerFn),
	}
	return srv
}
