package msg

import "fmt"

type Echo struct {
	message string
}

func (self *Echo) TypeID() uint16 {
	return uint16(1)
}

func (self *Echo) Encode() ([]byte, error) {
	return []byte(self.message), nil
}

func (self *Echo) Decode(bs []byte) error {
	self.message = string(bs)
	return nil
}

func NewEcho(obj interface{}) *Echo {
	return &Echo{message: fmt.Sprint(obj)}
}
