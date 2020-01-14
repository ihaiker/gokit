package remoting

import (
	"github.com/ihaiker/gokit/commons"
)

/*
	客户端连接管理器
*/
type ChannelManager interface {
	Add(channel Channel)
	Get(index interface{}) (channel Channel, has bool)
	Remove(channel Channel)
	Foreach(fn func(channel Channel))
}

type ipClientManager map[string]Channel

func NewIpClientManager() ChannelManager {
	return &ipClientManager{}
}

func (cm ipClientManager) Add(channel Channel) {
	cm[channel.GetRemoteAddress()] = channel
}

func (cm ipClientManager) Get(index interface{}) (channel Channel, has bool) {
	if commons.IsNil(index) {
		return
	}
	if ip, match := index.(string); match {
		channel, has = cm[ip]
	}
	return
}
func (cm ipClientManager) Remove(channel Channel) {
	ip := channel.GetRemoteAddress()
	delete(cm, ip)
	return
}

func (cm ipClientManager) Foreach(fn func(channel Channel)) {
	for _, c := range cm {
		fn(c)
	}
}
