package tcpKit

import (
    "fmt"
    "bufio"
    "io"
    "errors"
    "encoding/binary"
    "bytes"
    "reflect"
)

var (
    ErrInvalidProtocol = errors.New("invalid protocol")
)

type Protocol interface {
    Encode(msg interface{}) ([]byte, error)
    Decode(c io.Reader) (interface{}, error)
}

type ProtocolMaker func(c io.Reader) Protocol

type LineProtocol struct {
    reader    *bufio.Reader
    LineBreak string //换行分隔符
}

func (line *LineProtocol) Encode(msg interface{}) ([]byte, error) {
    return []byte(fmt.Sprintf("%s%s", msg, line.LineBreak)), nil
}

func (line *LineProtocol) Decode(c io.Reader) (interface{}, error) {
    if line.reader == nil {
        line.reader = bufio.NewReader(c)
    }
    ine, _, err := line.reader.ReadLine()
    return string(ine), err
}

type Package interface {
    ID() int16
    Encode() ([]byte, error)
    Decode(c io.Reader) (error)
}

type simpleProtocol struct {
    reg map[int16]reflect.Type
}

func NewSimpleProtocol() *simpleProtocol {
    return &simpleProtocol{reg: make(map[int16]reflect.Type)}
}

func (protocol *simpleProtocol) Reg(msg Package) {
    protocol.reg[msg.ID()] = reflect.TypeOf(msg)
}

func (protocol *simpleProtocol) Encode(msg interface{}) ([]byte, error) {
    if pkg, ok := msg.(Package); ok {
        if bs, err := pkg.Encode(); err != nil {
            return nil, err
        } else {
            bsWriter := new(bytes.Buffer)
            if err := binary.Write(bsWriter, binary.BigEndian, pkg.ID()); err != nil {
                return nil, err
            }
            if _, err := bsWriter.Write(bs); err != nil {
                return nil, err
            }
            return bsWriter.Bytes(), nil
        }
    }
    return nil, ErrInvalidProtocol
}

func (protocol *simpleProtocol) Decode(c io.Reader) (interface{}, error) {
    var id int16
    if err := binary.Read(c, binary.BigEndian, &id); err != nil {
        return nil, err
    }
    if pkgType, ok := protocol.reg[id]; ok {
        refType := reflect.New(pkgType.Elem())
        out := refType.MethodByName("Decode").Call([]reflect.Value{reflect.ValueOf(c)})
        if out[0].IsNil() {
            return refType.Interface(), nil
        } else {
            return nil, out[0].Interface().(error)
        }
        return refType.Interface(), nil
    }
    return nil, ErrInvalidProtocol
}
