package rpc

import (
	"github.com/ihaiker/gokit/net/buffer"
)

type Request struct {
	id      uint32
	URL     string
	Headers map[string]string
	Body    []byte
}

func (request *Request) ID() uint32 {
	return request.id
}

func (request *Request) TypeID() uint16 {
	return REQUEST
}

func (request *Request) Header(key, value string) (replace bool) {
	if request.Headers == nil {
		request.Headers = map[string]string{}
	}
	_, replace = request.Headers[key]
	request.Headers[key] = value
	return
}

func (request *Request) GetHeader(key string) (value string, has bool) {
	if request.Headers == nil {
		return
	}
	value, has = request.Headers[key]
	return
}

func (request *Request) Encode() ([]byte, error) {
	wr := buffer.NewWriter()
	_ = wr.UInt32(request.id)
	_ = wr.String(request.URL)
	if request.Headers == nil {
		_ = wr.UInt8(uint8(0))
	} else {
		_ = wr.UInt8(uint8(len(request.Headers)))
		for key, value := range request.Headers {
			_ = wr.String(key)
			_ = wr.String(value)
		}
	}
	if request.Body == nil {
		_ = wr.Write([]byte{})
	} else {
		_ = wr.Write(request.Body)
	}
	return wr.ToBytes(), nil
}

func (request *Request) Decode(bs []byte) (err error) {
	reader := buffer.NewReader(bs)
	if request.id, err = reader.UInt32(); err != nil {
		return
	}
	if request.URL, err = reader.String(); err != nil {
		return
	}
	headerLen := uint8(0)
	if headerLen, err = reader.UInt8(); err != nil {
		return
	}
	request.Headers = map[string]string{}
	if headerLen > 0 {
		for i := uint8(0); i < headerLen; i += 1 {
			var key, value string
			if key, err = reader.String(); err != nil {
				return
			}
			if value, err = reader.String(); err != nil {
				return
			}
			request.Headers[key] = value
		}
	}
	request.Body, err = reader.Bytes()
	return
}
