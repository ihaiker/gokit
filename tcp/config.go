package tcpKit

import "time"

type Config struct {
    PacketSendChanLimit    uint32 // the limit of packet send channel
    PacketReceiveChanLimit uint32 // the limit of packet receive channel
    AsynHandler            bool   // asynchronous handler message
    AcceptTimeout          time.Duration
    IdleTime, IdleTimeout  int //heartbeat time,and timeout 
}
