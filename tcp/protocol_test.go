package tcpKit

import (
    "testing"
    "github.com/docker/docker/pkg/testutil/assert"
    "bytes"
    "reflect"
    "io"
)

func TestRegisterProtocol(t *testing.T) {
    r := NewTVProtocol()
    show(t, r, nil)

    show(t, r, "[abcde]")
    show(t, r, []string{"A", "B"})

    show(t, r, io.EOF)
    show(t, r, []error{io.EOF, ErrWriteTimeout, ErrConnClosing})

    show(t, r, true)
    show(t, r, false)

    show(t, r, []bool{true, true, false, true})

    {
        a := int8(8)
        shows(t, r, a, &a, []int8{1, 2, 3})

        b := uint8(1)
        shows(t, r, b, &b, []uint8{1, 2, 3})
    }

    {
        a := int16(8)
        shows(t, r, a, &a, []int16{1, 2, 3})

        b := uint16(1)
        shows(t, r, b, &b, []uint16{1, 2, 3})
    }
}

func TestRegisterProtocols(t *testing.T) {
    r := NewTVProtocol()
    r.Reg(&ACK{})
    shows(t, r,
        nil,
        "[abcde]",
        []string{"A", "B"},

        io.EOF,
        []error{io.EOF, ErrWriteTimeout, ErrConnClosing},

        true,
        false,
        []bool{true, true, false, true},

        int8(8),
        []int8{4, -4},
        uint8(8),
        []uint8{4, 4},

        int16(16),
        []int16{8, -8},
        uint16(16),
        []uint16{8, 8},


        int32(16),
        uint32(16),

        int8(12),
        NewACK(1, 1), )
}

func TestACK(t *testing.T) {
    r := NewTVProtocol()
    ack := NewACK(1, 2)
    ack.Result = NewErrorACK(2, io.EOF)
    r.Reg(ack)
    show(t, r, ack)
}

func shows(t *testing.T, reg *regTVProtocol, msg ... interface{}) {
    size := len(msg)
    w := new(bytes.Buffer)
    for _, m := range msg {
        bs, err := reg.Encode(m)
        assert.NilError(t, err)
        w.Write(bs)
    }
    r := bytes.NewReader(w.Bytes())
    for i := 0; i < size; i++ {
        ret, err := reg.Decode(r)
        assert.NilError(t, err)
        t.Log(ret)
    }
}

func show(t *testing.T, reg *regTVProtocol, msg interface{}) {
    bs, err := reg.Encode(msg)
    assert.NilError(t, err)
    ret, err := reg.Decode(bytes.NewReader(bs))
    assert.NilError(t, err)
    t.Log(reflect.TypeOf(msg), msg, ret)
}

func TestName(t *testing.T) {

}
