package tlv

type Message interface {
	TypeID() uint16
	Encode() ([]byte, error)
	Decode([]byte) (error)
}
