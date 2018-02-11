package tcpKit

import (
    "io"
    "errors"
)

var (
    ErrInvalidProtocol = errors.New("invalid protocol")
)

type Protocol interface {
    Encode(msg interface{}) ([]byte, error)
    Decode(c io.Reader) (interface{}, error)
}

type ProtocolMaker func(c io.Reader) Protocol