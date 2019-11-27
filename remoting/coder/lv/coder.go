package lv

import (
	"encoding/binary"
	"errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/net/buffer"
	"github.com/ihaiker/gokit/remoting"
	"io"
)

var (
	ErrLengthMaxLimit = errors.New("LengthMaxLimit")
)

var logger = logs.GetLogger("lv")

type lv struct {
	//消息体最大长度
	maxLength uint16
}

func NewLVCoder(maxLength uint16) *lv {
	return &lv{maxLength: maxLength}
}

func (self *lv) Encode(channel remoting.Channel, msg interface{}) ([]byte, error) {
	message, match := msg.([]byte)
	if ! match {
		return nil, remoting.ErrInvalidArgument
	}
	length := uint16(len(message))
	if length > self.maxLength {
		return nil, ErrLengthMaxLimit
	}
	lengthBytes := buffer.UInt16(length)
	body := make([]byte, length+2)
	copy(body[0:2], lengthBytes)
	copy(body[2:], message)
	return body, nil
}

func (self *lv) Decode(channel remoting.Channel, c io.Reader) (interface{}, error) {
	var length uint16
	if err := binary.Read(c, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if length > self.maxLength {
		return nil, ErrLengthMaxLimit
	}

	bs := make([]byte, length)
	if _, err := io.ReadFull(c, bs); err != nil {
		return nil, err
	}
	return bs, nil
}
