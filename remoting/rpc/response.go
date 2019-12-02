package rpc

import (
	"errors"
	"github.com/ihaiker/gokit/net/buffer"
)

type Response struct {
	id    uint32
	Error error
	Body  []byte
}

func NewResponse(id uint32) *Response {
	return &Response{id: id}
}

func NewOKResponse(id uint32) *Response {
	resp := NewOKResponse(id)
	resp.Body = []byte("OK")
	return resp
}

func NewErrorResponse(id uint32, err error) *Response {
	resp := NewResponse(id)
	resp.Error = err
	return resp
}

func (response *Response) TypeID() uint16 {
	return RESPONSE
}

func (response *Response) Encode() ([]byte, error) {
	wr := buffer.NewWriter()
	_ = wr.UInt32(response.id)
	if response.Error == nil {
		_ = wr.String("")
	} else {
		_ = wr.String(response.Error.Error())
	}
	if response.Error == nil {
		if response.Body == nil {
			_ = wr.Write([]byte{})
		} else {
			_ = wr.Write(response.Body)
		}
	}
	return wr.ToBytes(), nil
}

func (response *Response) Decode(bs []byte) (err error) {
	reader := buffer.NewReader(bs)
	if response.id, err = reader.UInt32(); err != nil {
		return
	}
	var errString string
	if errString, err = reader.String(); err != nil {
		return
	}
	if errString != "" {
		response.Error = errors.New(errString)
	}
	if response.Error == nil {
		response.Body, err = reader.Bytes()
	}
	return
}
