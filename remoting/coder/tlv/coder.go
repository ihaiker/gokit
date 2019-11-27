package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"io"
	"reflect"
)

var (
	ErrLengthMaxLimit     = errors.New("LengthMaxLimit")
	ErrNotRegisterMessage = errors.New("NotRegisterMessage")
)

var logger = logs.GetLogger("tlv")

type tlvCoder struct {
	//消息体最大长度
	maxLength uint16
	reg       map[uint16]reflect.Type
}

func NewTLVCoder(maxLength uint16) *tlvCoder {
	return &tlvCoder{maxLength: maxLength, reg: make(map[uint16]reflect.Type)}
}

func (self *tlvCoder) Reg(msg Message) error {
	if msg == nil {
		return remoting.ErrInvalidArgument
	}
	self.reg[msg.TypeID()] = reflect.TypeOf(msg)
	return nil
}

func (self *tlvCoder) Encode(channel remoting.Channel, msg interface{}) ([]byte, error) {
	message, match := msg.(Message)
	if ! match {
		return nil, remoting.ErrInvalidArgument
	}

	typeId := message.TypeID()
	if _, has := self.reg[typeId]; !has {
		return nil, ErrNotRegisterMessage
	}

	w := new(bytes.Buffer)

	if bs, err := message.Encode(); err != nil {
		return nil, err
	} else {
		length := uint16(len(bs))
		if length > self.maxLength {
			return nil, ErrLengthMaxLimit
		}

		//type
		if err = binary.Write(w, binary.BigEndian, typeId); err != nil {
			return nil, err
		}
		//length
		if err = binary.Write(w, binary.BigEndian, length); err != nil {
			return nil, err
		}

		//value
		if length > 0 {
			if err = binary.Write(w, binary.BigEndian, bs); err != nil {
				return nil, err
			}
		}
		w.Cap()
		return w.Bytes(), nil
	}
}

func (self *tlvCoder) Decode(channel remoting.Channel, c io.Reader) (interface{}, error) {
	var typeId uint16 = 0
	if err := binary.Read(c, binary.BigEndian, &typeId); err != nil {
		return nil, err
	}

	msgType, has := self.reg[typeId]
	if ! has {
		return nil, ErrNotRegisterMessage
	}

	var length uint16
	if err := binary.Read(c, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	if length > self.maxLength {
		return nil, ErrLengthMaxLimit
	}

	bs := make([]byte, length)
	if length > 0 {
		if _, err := io.ReadFull(c, bs); err != nil {
			return nil, err
		}
	}
	
	var msgValue reflect.Value
	if msgType.Kind() == reflect.Ptr {
		msgValue = reflect.New(msgType.Elem())
	} else {
		msgValue = reflect.New(msgType)
	}
	out := msgValue.MethodByName("Decode").Call([]reflect.Value{reflect.ValueOf(bs)})

	if out[0].IsNil() {
		return msgValue.Interface(), nil
	} else {
		return nil, out[0].Interface().(error)
	}
}
