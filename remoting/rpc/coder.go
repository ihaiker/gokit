package rpc

import (
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/tlv"
)

func coderMaker(channel remoting.Channel) remoting.Coder {
	return newCoder()
}

func newCoder() remoting.Coder {
	coder := tlv.NewTLVCoder(1024)
	_ = coder.Reg(new(Ping))
	_ = coder.Reg(new(Pong))
	_ = coder.Reg(new(Request))
	_ = coder.Reg(new(Response))
	return coder
}
