package tcpKit

import (
    "io"
    "reflect"
    "bytes"
    "encoding/binary"
    "github.com/ihaiker/gokit/commons"
    "errors"
)

type Package interface {
    PID() int16
    Encode() ([]byte, error)
    Decode(c io.Reader) (error)
}

type regTVProtocol struct {
    reg map[int16]reflect.Type
}

func NewTVProtocol() *regTVProtocol {
    return &regTVProtocol{reg: make(map[int16]reflect.Type)}
}

func (protocol *regTVProtocol) Reg(msg Package) error {
    if msg == nil || msg.PID() <= REGISTER_TYPE_MAX {
        return ErrInvalidProtocol
    } else {
        protocol.reg[msg.PID()] = reflect.TypeOf(msg)
        return nil
    }
}

func (protocol *regTVProtocol) writeType(bsWriter *bytes.Buffer, t int16) {
    err := binary.Write(bsWriter, binary.BigEndian, t)
    commonKit.IfPanic(err)
}

func (protocol *regTVProtocol) Encode(msg interface{}) (bs []byte, err error) {
    defer func() { err = commonKit.DCatch(recover(), err) }()

    w := new(bytes.Buffer)

    if msg == nil {
        protocol.writeType(w, REGISTER_TYPE_NIL)
        return w.Bytes(), nil
    }

    switch msg.(type) {
    case error:
        protocol.writeType(w, REGISTER_TYPE_ERROR)
        WriteString(w, msg.(error).Error())
    case []error:
        protocol.writeType(w, REGISTER_TYPE_ERROR_ARRAY)
        ary := msg.([]error)
        binary.Write(w, binary.BigEndian, uint8(len(ary)))
        for _, s := range ary {
            WriteString(w, s.Error())
        }

    case string:
        protocol.writeType(w, REGISTER_TYPE_STRING)
        WriteString(w, msg.(string))
    case []string:
        protocol.writeType(w, REGISTER_TYPE_STRING_ARRAY)
        ary := msg.([]string)
        binary.Write(w, binary.BigEndian, uint8(len(ary)))
        for _, s := range ary {
            WriteString(w, s)
        }

    case bool, *bool:
        protocol.writeType(w, REGISTER_TYPE_BOOL)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []bool:
        protocol.writeType(w, REGISTER_TYPE_BOOL_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]bool)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case int8, *int8:
        protocol.writeType(w, REGISTER_TYPE_INT8)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []int8:
        protocol.writeType(w, REGISTER_TYPE_INT8_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]int8)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case uint8, *uint8:
        protocol.writeType(w, REGISTER_TYPE_UINT8)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []uint8:
        protocol.writeType(w, REGISTER_TYPE_UINT8_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]uint8)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case int16, *int16:
        protocol.writeType(w, REGISTER_TYPE_INT16)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []int16:
        protocol.writeType(w, REGISTER_TYPE_INT16_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]int16)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case uint16, *uint16:
        protocol.writeType(w, REGISTER_TYPE_UINT16)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []uint16:
        protocol.writeType(w, REGISTER_TYPE_UINT16_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]uint16)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case int32, *int32:
        protocol.writeType(w, REGISTER_TYPE_INT32)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []int32:
        protocol.writeType(w, REGISTER_TYPE_INT32_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]int32)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case uint32, *uint32:
        protocol.writeType(w, REGISTER_TYPE_UINT32)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []uint32:
        protocol.writeType(w, REGISTER_TYPE_UINT32_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]uint32)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case int64, *int64:
        protocol.writeType(w, REGISTER_TYPE_INT64)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []int64:
        protocol.writeType(w, REGISTER_TYPE_INT64_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]int64)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case uint64, *uint64:
        protocol.writeType(w, REGISTER_TYPE_UINT64)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []uint64:
        protocol.writeType(w, REGISTER_TYPE_UINT64_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]uint64)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case int, *int:
        protocol.writeType(w, REGISTER_TYPE_INT)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []int:
        protocol.writeType(w, REGISTER_TYPE_INT_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]int)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    case uint, *uint:
        protocol.writeType(w, REGISTER_TYPE_UINT)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))
    case []uint:
        protocol.writeType(w, REGISTER_TYPE_UINT_ARRAY)
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, uint8(len(msg.([]uint)))))
        commonKit.IfPanic(binary.Write(w, binary.BigEndian, msg))

    default:
        if pkg, ok := msg.(Package); ok {
            if bs, err := pkg.Encode(); err != nil {
                return nil, err
            } else {
                if err := binary.Write(w, binary.BigEndian, pkg.PID()); err != nil {
                    return nil, err
                }
                if _, err := w.Write(bs); err != nil {
                    return nil, err
                }
                return w.Bytes(), nil
            }
        }
        return nil, ErrInvalidProtocol
    }
    return w.Bytes(), nil
}

func (protocol *regTVProtocol) Decode(c io.Reader) (interface{}, error) {
    var typeId int16
    if err := binary.Read(c, binary.BigEndian, &typeId); err != nil {
        return nil, err
    }
    switch typeId {
    case REGISTER_TYPE_NIL:
        return nil, nil

    case REGISTER_TYPE_STRING:
        return ReadString(c), nil
    case REGISTER_TYPE_STRING_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        ary := make([]string, size)
        for i := 0; i < int(size); i++ {
            ary[i] = ReadString(c)
        }
        return ary, nil

    case REGISTER_TYPE_ERROR:
        return errors.New(ReadString(c)), nil
    case REGISTER_TYPE_ERROR_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        ary := make([]error, size)
        for i := 0; i < int(size); i++ {
            ary[i] = errors.New(ReadString(c))
        }
        return ary, nil

    case REGISTER_TYPE_BOOL:
        var p bool
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_BOOL_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]bool, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_INT8:
        var p int8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_INT8_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]int8, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_UINT8:
        var p uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_UINT8_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]uint8, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_INT16:
        var p int16
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_INT16_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]int16, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_UINT16:
        var p uint16
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_UINT16_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]uint16, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_INT32:
        var p int32
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_INT32_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]int32, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_UINT32:
        var p uint32
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_UINT32_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]uint32, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_INT64:
        var p int64
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_INT64_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]int64, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_UINT64:
        var p uint64
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_UINT64_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]uint64, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_INT:
        var p int
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_INT_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]int, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_UINT:
        var p uint
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil
    case REGISTER_TYPE_UINT_ARRAY:
        var size uint8
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &size))
        p := make([]uint, size)
        commonKit.IfPanic(binary.Read(c, binary.BigEndian, &p))
        return p, nil

    case REGISTER_TYPE_MAX:
        return nil, ErrInvalidProtocol
    default:
        if pkgType, ok := protocol.reg[typeId]; ok {
            refType := reflect.New(pkgType.Elem())
            out := refType.MethodByName("Decode").Call([]reflect.Value{reflect.ValueOf(c)})
            if out[0].IsNil() {
                return refType.Interface(), nil
            } else {
                return nil, out[0].Interface().(error)
            }
            return refType.Interface(), nil
        }
    }
    return nil, ErrInvalidProtocol
}

func ReadString(conn io.Reader) string {
    var len uint8
    commonKit.IfPanic(binary.Read(conn, binary.BigEndian, &len))
    if len == 0 {
        return ""
    }
    val := make([]byte, len)
    if redLen, err := io.ReadFull(conn, val); err != nil {
        panic(err)
    } else if uint8(redLen) != len {
        panic(errors.New("消息长度不足："))
    }
    return string(val)
}

func WriteString(w io.Writer, str string) {
    var len uint8 = uint8(len(str))
    commonKit.IfPanic(binary.Write(w, binary.BigEndian, len))
    if len > 0 {
        _, err := w.Write([]byte(str))
        commonKit.IfPanic(err)
    }
}
