package remoting

import (
	"io"
)

type Coder interface {
	Encode(channel Channel, msg interface{}) ([]byte, error)
	Decode(channel Channel, reader io.Reader) (interface{}, error)
}

type CoderMaker func(channel Channel) Coder
