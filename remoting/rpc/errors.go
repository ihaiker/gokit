package rpc

import (
	"errors"
)

var (
	ErrNotFount = errors.New("not found")
	ErrSystemError = errors.New("SystemError")
	ErrRpcTimeout = errors.New("rpc timeout")
)

