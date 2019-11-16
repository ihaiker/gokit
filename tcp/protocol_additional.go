package tcpKit

import (
    "fmt"
    "bytes"
    "github.com/ihaiker/gokit/commons"
    "io"
    "encoding/binary"
    "errors"
)

const (
    REGISTER_TYPE_ACK  = REGISTER_TYPE_MAX + 1 + iota
    REGISTER_TYPE_PING
    REGISTER_TYPE_PONG
)

type Idle uint16

func (idle Idle) String() string {
    if uint16(idle) == REGISTER_TYPE_PING {
        return "PING"
    }
    return "PONG"
}

func (idle Idle) PID() uint16 {
    return uint16(idle)
}

func (idle Idle) Encode(protocol Protocol) ([]byte, error) {
    return []byte{}, nil
}

func (idle Idle) Decode(protocol Protocol, c io.Reader) (error) {
    return nil
}

const (
    PING Idle = Idle(REGISTER_TYPE_PING)
    PONG Idle = Idle(REGISTER_TYPE_PONG)
)

type ACK struct {
    SendId int64
    Err    error
    Result interface{}
}

func (ack *ACK) String() string {
    return fmt.Sprintf("ACK:{SendId: %d, Err: %v, Result: %v}", ack.SendId, ack.Err, ack.Result)
}

func (ack *ACK) PID() uint16 {
    return REGISTER_TYPE_ACK
}

func (ack *ACK) Encode(protocol Protocol) ([]byte, error) {
    buffer := new(bytes.Buffer)
    commons.PanicIfPresent(binary.Write(buffer, binary.BigEndian, ack.SendId))

    //error
    if ack.Err == nil {
        WriteString(buffer, "")
    } else {
        WriteString(buffer, ack.Err.Error())
    }

    //result
    bs, err := protocol.Encode(ack.Result)
    commons.PanicIfPresent(err)
    _, err = buffer.Write(bs)
    commons.PanicIfPresent(err)

    return buffer.Bytes(), nil
}
func (ack *ACK) Decode(protocol Protocol, c io.Reader) (error) {
    if err := binary.Read(c, binary.BigEndian, &ack.SendId); err != nil {
        return err
    }
    eStr := ReadString(c)
    if eStr == "" {

    } else {
        ack.Err = errors.New(eStr)
    }

    if ret, err := protocol.Decode(c); err != nil {
        return err
    } else if ret != nil {
        ack.Result = ret
    }
    return nil
}

func NewACK(sendId int64, ret interface{}) *ACK {
    return &ACK{SendId: sendId, Result: ret}
}
func NewErrorACK(sendId int64, err error) *ACK {
    return &ACK{SendId: sendId, Err: err}
}

//
//用户认证消息
type IDWapper struct {
    SendId int64 //认证消息发送ID
}

func (self *IDWapper) Encode(p Protocol) ([]byte, error) {
    w := new(bytes.Buffer)
    commons.PanicIfPresent(binary.Write(w, binary.BigEndian, self.SendId))
    if err := self.EncodeEntry(p, w); err != nil {
        return nil, err
    }
    return w.Bytes(), nil
}

func (self *IDWapper) EncodeEntry(p Protocol, buf *bytes.Buffer) error {
    return nil
}

func (evt *IDWapper) Decode(p Protocol, c io.Reader) error {
    commons.PanicIfPresent(binary.Read(c, binary.BigEndian, &evt.SendId))
    return evt.DecodeEntry(p, c)
}

func (evt *IDWapper) DecodeEntry(p Protocol, c io.Reader) error {
    return nil
}
