package remoting

import (
	"encoding/json"
)

type Options struct {
	//异步发送chan定义长度
	SendChanLimit int `json:"sendChanLimit" yaml:"sendChanLimit" yaml:"sendChanLimit"`

	//接收chan容量
	ReceiveChanLimit int `json:"receiveChanLimit" yaml:"receiveChanLimit" yaml:"receiveChanLimit"`

	//是否异步处理接收的消息，队列大小
	WorkerGroup int `json:"workerGroup" yaml:"workerGroup" toml:"workerGroup"`

	//心跳检测周期,单位秒
	IdleTimeSeconds int `json:"idleTimeSeconds" yaml:"idleTimeSeconds" toml:"idleTimeSeconds"`
	//心跳检测超时次数，多少次检测后失效
	IdleTimeout int `json:"idleTimeout" yaml:"idleTimeout" toml:"idleTimeout"`

	SendBuf int `json:"sendBuf" yaml:"sendBuf" toml:"sendBuf"`
	RecvBuf int `json:"recvBuf" yaml:"recvBuf" toml:"recvBuf"`
}

func (cfg *Options) String() string {
	bs, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func DefaultOptions() *Options {
	return &Options{
		SendChanLimit: 1000, ReceiveChanLimit: 1000,
		WorkerGroup:     100,
		IdleTimeSeconds: 3, IdleTimeout: 3,
		SendBuf: 1024, RecvBuf: 1024,
	}
}
