package remoting

type Handler interface {
	//链接事件，当客户端链接时调用
	OnConnect(c Channel)

	//新消息事件，当客户端发来新的消息时调用
	OnMessage(s Channel, msg interface{})

	//编码异常事件，当编码消息时
	OnEncodeError(s Channel, msg interface{}, err error)

	//处理错误事件，当OnMessage抛出未能处理的错误
	OnError(s Channel, msg interface{}, err error)

	//解码异常事件，解码时错误
	OnDecodeError(s Channel, err error)

	//发送心跳包
	OnIdle(c Channel)

	//关闭事件，当当前客户端关闭连接
	OnClose(c Channel)
}

type HandlerMaker func(channel Channel) Handler
