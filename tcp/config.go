package tcpKit

import (
    "encoding/json"
    "time"
)

type Config struct {
    // the limit of packet send channel
    PacketSendChanLimit uint32 `json:"packetSendChanLimit" yaml:"packetSendChanLimit" toml:"packetSendChanLimit"`

    // the limit of packet receive channel
    PacketReceiveChanLimit uint32 `json:"packetReceiveChanLimit" yaml:"packetReceiveChanLimit" toml:"packetReceiveChanLimit"`

    // asynchronous handler message
    AsynHandler   bool          `json:"asynHandler" yaml:"asynHandler" toml:"asynHandler"`
    AcceptTimeout time.Duration `json:"acceptTimeout" yaml:"acceptTimeout" toml:"acceptTimeout"`

    //heartbeat time,and timeout
    IdleTime    int `json:"idleTime" yaml:"idleTime" toml:"idleTime"`
    IdleTimeout int `json:"idleTimeout" yaml:"idleTimeout" toml:"idleTimeout"`
}

func (cfg *Config) String() string {
    bs, _ := json.Marshal(cfg)
    return string(bs)
}

func DefaultTCPConfig() *Config {
    return &Config{
        PacketReceiveChanLimit: 10,
        PacketSendChanLimit:    10,
        AcceptTimeout:          time.Second,
        IdleTime:               1000,
        IdleTimeout:            3,
    }
}
