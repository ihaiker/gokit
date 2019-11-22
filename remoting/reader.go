package remoting

import "net"

type ReaderBuffer struct {
	connect *net.TCPConn
}

func (r *ReaderBuffer) Read(p []byte) (int, error) {
	return r.connect.Read(p)
}
