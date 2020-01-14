package remoting

type Event struct {
	Type    EventType
	Channel Channel
	Values  []interface{}
}

type EventType string

const (
	ConnectEvent   EventType = "ConnectEvent"
	CloseEvent     EventType = "CloseEvent"
	DecodeErrEvent EventType = "DecodeErrEvent"
	EncodeErrEvent EventType = "EncodeErrEvent"
	ErrEvent       EventType = "ErrEvent"
	MessageEvent   EventType = "MessageEvent"
	IdleEvent      EventType = "IdleEvent"
)

func (event EventType) String() string {
	return string(event)
}

func NewEvent(name EventType, ch Channel, values ...interface{}) *Event {
	return &Event{
		Type: name, Channel: ch, Values: values,
	}
}

type Handler interface {
	OnEvent(event *Event)
}

type HandlerMaker func(channel Channel) Handler
