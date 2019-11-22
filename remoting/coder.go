package remoting

import (
	"errors"
	"io"
)

var (
	ErrInvalidCoder = errors.New("Err: invalid coder")
)

type Coder interface {
	Encode(channel Channel, msg interface{}) ([]byte, error)
	Decode(channel Channel, reader io.Reader) (interface{}, error)
}

type CoderMaker func(channel Channel) Coder
