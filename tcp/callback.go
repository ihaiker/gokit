package gotcp

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface {
	// OnConnect is called when the connection was accepted,
	OnConnect(*Conn)

	// OnMessage is called when the connection receives a packet
	OnMessage(*Conn, interface{})

	OnIdle(*Conn, IdleState)

	// OnError is called where the connection handle the received packet
	OnError(*Conn, interface{})

	// OnClose is called when the connection closed
	OnClose(*Conn)
}

type DefConnCallback struct{}

func (self *DefConnCallback) OnConnect(conn *Conn) {}
func (self *DefConnCallback) OnMessage(conn *Conn, msg interface{}) {}
func (self *DefConnCallback) OnIdle(conn *Conn, idleState IdleState) {}
func (self *DefConnCallback) OnError(conn *Conn, err interface{}) {}
func (self *DefConnCallback) OnClose(conn *Conn) {}

type ClientCallback interface {
	// OnConnect is called when the connection was accepted,
	OnConnect(*Client)

	// OnMessage is called when the connection receives a packet
	OnMessage(*Client, interface{})

	// OnError is called where the connection handle the received packet
	OnError(*Client, interface{})

	// OnClose is called when the connection closed
	OnClose(*Client)
}
type DefClientCallback struct{}

func (self *DefClientCallback) OnConnect(client *Client) {}
func (self *DefClientCallback) OnMessage(*Client, interface{}) {}
func (self *DefClientCallback) OnError(*Client, interface{}) {}
func (self *DefClientCallback) OnClose(*Client) {}