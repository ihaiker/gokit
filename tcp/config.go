package gotcp

const (
	PACKET_SEND_CHAN_LIMIT = 20
	PACKET_RECEIVE_CHAN_LIMIT = 20
)

type Config struct {
	PacketSendChanLimit    uint32 // the limit of packet send channel
	PacketReceiveChanLimit uint32 // the limit of packet receive channel
	AsyncMessageHand       bool   // asynchronous message handling
}

func DefConfig() *Config {
	cfg := &Config{
		PacketSendChanLimit     : PACKET_SEND_CHAN_LIMIT,
		PacketReceiveChanLimit  : PACKET_RECEIVE_CHAN_LIMIT,
		AsyncMessageHand        : false,
	}
	return cfg
}