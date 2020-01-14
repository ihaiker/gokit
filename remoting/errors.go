package remoting

import "errors"

var (
	ErrConnectClosed  = errors.New("use of closed network connection")
	ErrWriteTimeout   = errors.New("write operation timed out")
	ErrInvalidMessage = errors.New("Invalid message")
)
