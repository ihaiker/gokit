package msg

import (
	"github.com/ihaiker/gokit/logs"
	"io"
)

type Package2 struct {
    Msg string
}

func (p *Package2) PID() uint16 {
    return 1
}
func (p *Package2) Encode() ([]byte, error) {
    return ([]byte(p.Msg))[0:4], nil
}
func (p *Package2) Decode(c io.Reader) (error) {
    logs.Info("Decode Info")
    bs := make([]byte, 4)
    _, err := io.ReadFull(c, bs)
    p.Msg = string(bs)
    return err
}
