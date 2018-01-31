package tcpKit

type simpleWrapper struct {
    handler  *regHandler
    protocol *regTVProtocol
}

func NewRegister() *simpleWrapper {
    r := &simpleWrapper{}
    r.handler = newRegHandler()
    r.protocol = NewTVProtocol()
    return r
}

func (w *simpleWrapper) Reg(msg Package) {
    w.protocol.Reg(msg)
}

func (w *simpleWrapper) OnConnect(OnConnect func(c *Connect)) {
    w.handler.onConnect = OnConnect
}
func (w *simpleWrapper) OnMessage(OnMessage func(c *Connect, msg interface{})) {
    w.handler.onMessage = OnMessage
}
func (w *simpleWrapper) OnEncodeError(OnEncodeError func(c *Connect, msg interface{}, err error)) {
    w.handler.onEncodeError = OnEncodeError
}
func (w *simpleWrapper) OnError(OnError func(c *Connect, err error, msg interface{})) {
    w.handler.onError = OnError
}
func (w *simpleWrapper) OnDecodeError(OnDecodeError func(c *Connect, err error)) {
    w.handler.onDecodeError = OnDecodeError
}
func (w *simpleWrapper) OnIdle(OnIdle func(c *Connect)) {
    w.handler.onIdle = OnIdle
}
func (w *simpleWrapper) OnClose(OnClose func(c *Connect)) {
    w.handler.onClose = OnClose
}
