package tcpKit

import (
    "fmt"
    "bytes"
    "github.com/ihaiker/gokit/commons"
    "io"
    "encoding/binary"
)

const (
    REGISTER_TYPE_ACK = REGISTER_TYPE_MAX + 1 + iota
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

func (idle Idle) Decode(protocol Protocol,c io.Reader) (error) {
    return nil
}

const (
    PING Idle = Idle(REGISTER_TYPE_PING)
    PONG Idle = Idle(REGISTER_TYPE_PONG)
)

type ACK struct {
    SendId, ACKID int64
    Err           error
    Result        interface{}
}

func (ack *ACK) String() string {
    return fmt.Sprintf("ACK:{SendId: %d, ACKID: %d, Err: %v, Result: %v}", ack.SendId, ack.ACKID, ack.Err, ack.Result)
}

func (ack *ACK) PID() uint16 {
    return REGISTER_TYPE_ACK
}

func (ack *ACK) Encode(protocol Protocol) ([]byte, error) {
    buffer := new(bytes.Buffer)
    commonKit.IfPanic(binary.Write(buffer, binary.BigEndian, ack.SendId))
    commonKit.IfPanic(binary.Write(buffer, binary.BigEndian, ack.ACKID))

    //error
    bs, err := protocol.Encode(ack.Err)
    commonKit.IfPanic(err)
    _, err = buffer.Write(bs)
    commonKit.IfPanic(err)

    //result
    bs, err = protocol.Encode(ack.Result)
    commonKit.IfPanic(err)
    _, err = buffer.Write(bs)
    commonKit.IfPanic(err)

    return buffer.Bytes(), nil
}
func (ack *ACK) Decode(protocol Protocol, c io.Reader) (error) {
    if err := binary.Read(c, binary.BigEndian, &ack.SendId); err != nil {
        return err
    }
    if err := binary.Read(c, binary.BigEndian, &ack.ACKID); err != nil {
        return err
    }
    if errOut, err := protocol.Decode(c); err != nil {
        return err
    } else if errOut != nil {
        ack.Err = errOut.(error)
    }
    if ret, err := protocol.Decode(c); err != nil {
        return err
    } else if ret != nil {
        ack.Result = ret
    }
    return nil
}

func NewACK(sendId, ackID int64, errMsg error) *ACK {
    return &ACK{SendId: sendId, ACKID: ackID, Err: errMsg}
}
