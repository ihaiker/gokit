package gotcp

import "errors"


// Error type
var (
	ERR_CONN_CLOSING = errors.New("use of closed network connection")
	ERR_WRITE_BLOCKING = errors.New("write packet was blocking")
	ERR_READ_BLOCKING = errors.New("read packet was blocking")
)

type DecodePackageError struct {
	Msg interface{}
}
type EncodePackageError struct {
	Msg interface{}
}