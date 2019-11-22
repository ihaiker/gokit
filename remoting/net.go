package remoting

import (
	"net"
	"strings"
)

func DialTcp(address string) (net.Conn, error) {
	network := "tcp4"
	ip := address
	if idx := strings.Index(address, "://"); idx != -1 {
		network = address[:idx]
		ip = address[idx+3:]
	}
	if _, err := net.ResolveTCPAddr(network, ip); err != nil {
		return nil, err
	} else {
		return net.Dial(network, address)
	}
}

func DialUnit(address string) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", address)
	if err != nil {
		return nil, err
	}
	return net.DialUnix("unix", nil, unixAddr)
}

func Dial(address string) (net.Conn, error) {
	if strings.HasPrefix(address, "tcp://") ||
		strings.HasPrefix(address, "tcp4://") ||
		strings.HasPrefix(address, "tcp6://") {
		return DialTcp(address)
	} else if strings.HasPrefix(address, "unix://") {
		return DialUnit(address[6:])
	}
	return DialTcp(address)
}

func ListenTcp(address string) (*net.TCPListener, error) {
	network := "tcp4"
	ip := address
	if idx := strings.Index(address, "://"); idx != -1 {
		network = address[:idx]
		ip = address[idx+3:]
	}
	if tcpAddr, err := net.ResolveTCPAddr(network, ip); err != nil {
		return nil, err
	} else {
		return net.ListenTCP("tcp", tcpAddr);
	}
}

func ListenUnix(address string) (net.Listener, error) {
	add, err := net.ResolveUnixAddr("unix", address)
	if err != nil {
		return nil, err
	}
	return net.ListenUnix("unix", add)
}

func Listen(address string) (net.Listener, error) {
	if strings.HasPrefix(address, "tcp://") ||
		strings.HasPrefix(address, "tcp4://") ||
		strings.HasPrefix(address, "tcp6://") {
		return ListenTcp(address)
	} else if strings.HasPrefix(address, "unix://") {
		return ListenUnix(address[6:])
	}
	return ListenTcp(address)
}
