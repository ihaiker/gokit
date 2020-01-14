package remoting

import "io"

type InboundCoder interface {
	Decode(channel Channel, reader io.Reader) (Message, error)
}

type OutboundCoder interface {
	Encode(channel Channel, msg Message) ([]byte, error)
}

type Coder interface {
	InboundCoder
	OutboundCoder
}

type CoderMaker func(channel Channel) Coder
