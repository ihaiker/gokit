package main

import (
	"net"
	"github.com/ihaiker/gokit/tcp/examples/echo"
	"github.com/ihaiker/gokit/tcp"
	"fmt"
	"log"
	"time"
)

func checkClientError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type EchoClientCallback struct {
	gotcp.ClientCallback
}

func (self *EchoClientCallback) OnConnect(client *gotcp.Client) {
	fmt.Println(client.GetRawConn().RemoteAddr(), "connect !")
	client.Write(echo.NewEchoPacket([]byte("hello")))
}

func (self *EchoClientCallback) OnMessage(client *gotcp.Client, p interface{}) {
	echoPack := p.(*echo.EchoPacket)
	fmt.Println("client receiver : ", string(echoPack.GetBody()))
}

func (self *EchoClientCallback) OnClose(client *gotcp.Client) {
	fmt.Println("the client close!")
	time.Sleep(time.Second*3)
}

func main() {
	tcp_address, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkClientError(err)
	conn, err := net.DialTCP("tcp", nil, tcp_address)
	checkClientError(err)

	conn.SetKeepAlive(true)
	conn.SetNoDelay(true)

	protocol := &echo.EchoProtocol{}
	callback := &EchoClientCallback{}
	config := gotcp.DefConfig()

	client := gotcp.NewClient(config, callback, protocol)
	go client.Start(conn)

	time.Sleep(time.Second)

	client.Close()

}
