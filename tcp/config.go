package tcpKit

import "time"

type Config struct {
    // the limit of packet send channel
    PacketSendChanLimit uint32 `json:"packet_send_chan_limit" yaml:"packet_send_chan_limit"`
    // the limit of packet receive channel
    PacketReceiveChanLimit uint32 `json:"packet_receive_chan_limit" yaml:"packet_receive_chan_limit"`
    // asynchronous handler message
    AsynHandler   bool          `json:"asyn_handler" yaml:"asyn_handler"`
    AcceptTimeout time.Duration `json:"accept_timeout" yaml:"accept_timeout"`
    //heartbeat time,and timeout
    IdleTime    int `json:"idle_time" yaml:"idle_time"`
    IdleTimeout int `json:"idle_timeout" yaml:"idle_timeout"`
}
