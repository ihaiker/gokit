package remoting

import (
	"net"
	"strings"
)

func Dial(address string) (net.Conn, error) {
	if strings.HasPrefix(address, "tcp://") {
		return DialTcp(address)
	} else if strings.HasPrefix(address, "unix://") {
		return DialUnix(address)
	}
	return DialTcp(address)
}

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

func DialUnix(address string) (net.Conn, error) {
	if strings.HasPrefix(address, "unix://") {
		address = address[6:]
	}
	unixAddr, err := net.ResolveUnixAddr("unix", address)
	if err != nil {
		return nil, err
	}
	return net.DialUnix("unix", nil, unixAddr)
}

func Listen(address string) (net.Listener, error) {
	if strings.HasPrefix(address, "tcp://") {
		return ListenTcp(address)
	} else if strings.HasPrefix(address, "unix://") {
		return ListenUnix(address)
	}
	return ListenTcp(address)
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
		return net.ListenTCP("tcp", tcpAddr)
	}
}

func ListenUnix(address string) (net.Listener, error) {
	if strings.HasPrefix(address, "unix://") {
		address = address[6:]
	}
	add, err := net.ResolveUnixAddr("unix", address)
	if err != nil {
		return nil, err
	}
	return net.ListenUnix("unix", add)
}
