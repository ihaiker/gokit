package remoting

import (
	"encoding/json"
)

type Config struct {
	//异步发送chan定义长度
	SendChanLimit int `json:"sendChanLimit" yaml:"sendChanLimit" yaml:"sendChanLimit"`

	//接收chan容量
	ReceiveChanLimit int `json:"receiveChanLimit" yaml:"receiveChanLimit" yaml:"receiveChanLimit"`

	//是否异步处理接收的消息，队列大小
	AsyncHandlerGroup int `json:"asyncHandlerGroup" yaml:"asyncHandlerGroup" toml:"asyncHandlerGroup"`

	//心跳检测周期,单位秒
	IdleDuration int `json:"idleDuration" yaml:"idleDuration" toml:"idleDuration"`
	//心跳检测超时次数，多少次检测后失效
	IdleTimeout int `json:"idleTimeout" yaml:"idleTimeout" toml:"idleTimeout"`

	//读缓存大小
	ReadBufferSize int `json:"readBufferSize" yaml:"readBufferSize" toml:"readBufferSize"`

	//写缓存大小
	WriteBufferSize int `json:"writeBufferSize" yaml:"writeBufferSize" yaml:"writeBufferSize"`
}

func (cfg *Config) String() string {
	bs, _ := json.Marshal(cfg)
	return string(bs)
}

func DefaultTCPConfig() *Config {
	return &Config{
		SendChanLimit: 1000, ReceiveChanLimit: 1000,
		AsyncHandlerGroup: 100,
		IdleDuration:      3,
		IdleTimeout:       3,
		WriteBufferSize:   512,
		ReadBufferSize:    512,
	}
}
