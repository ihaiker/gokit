package gotcp

import (
	"net"
)


type Protocol interface {
	Decode(conn *net.TCPConn) (interface{}, error)
	Encode(pkg interface{}) ([]byte, error)
}
