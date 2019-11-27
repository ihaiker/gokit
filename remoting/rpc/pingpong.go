package rpc

type Ping struct {
}

func (request *Ping) TypeID() uint16 {
	return PING
}

func (request *Ping) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (request *Ping) Decode(bs []byte) (err error) {
	return nil
}

type Pong struct {
}

func (request *Pong) TypeID() uint16 {
	return PONG
}

func (request *Pong) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (request *Pong) Decode(bs []byte) (err error) {
	return nil
}
