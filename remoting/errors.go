package remoting

import "errors"

var (
	ErrConnectClosed   = errors.New("use of closed network connection")
	ErrWriteTimeout    = errors.New("write operation timed out")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidHandler  = errors.New("invalid handler")
	ErrInvalidCoder    = errors.New("invalid coder")
	ErrNoServerConnect = errors.New("no server to connect")
)
