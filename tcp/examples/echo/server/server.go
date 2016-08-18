package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	tcpkit "github.com/ihaiker/gokit/tcp"
	"github.com/ihaiker/gokit/tcp/examples/echo"
	ktime "github.com/ihaiker/gokit/time"
)



type Callback struct{
	tcpkit.DefConnCallback
}

func (this *Callback) OnConnect(c *tcpkit.Conn) {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
}

func (this *Callback) OnMessage(c *tcpkit.Conn, p interface{}) {
	echoPacket := p.(*echo.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.Write(echo.NewEchoPacket(echoPacket.GetBody()), time.Second)
}

func (this *Callback) OnError(c *tcpkit.Conn, err interface{}) {
	fmt.Println(err)
}

func (this *Callback) OnClose(c *tcpkit.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}


type EchoHeartbeatHandler struct{}

func (self *EchoHeartbeatHandler) OnIdle(idleState tcpkit.IdleState, conn *tcpkit.Conn) {
	switch idleState {
	case tcpkit.READ_IDLE_STATE:
		t := ktime.JavaFormat(time.Now(), "yyyy/MM/dd HH:mm:ss")
		fmt.Println("read idle status", t)
		p := echo.NewEchoPacket([]byte(t))
		conn.Write(p, time.Second)

	case tcpkit.WRITER_IDLE_STATE:
		fmt.Println("write idle status", ktime.JavaFormat(time.Now(), "yyyy/MM/dd HH:mm:ss"))

	case tcpkit.ALL_IDLE_STATE:
		fmt.Println("all idle status", ktime.JavaFormat(time.Now(), "yyyy/MM/dd HH:mm:ss"))
	}
}

func hb_maker(c *net.TCPConn) *tcpkit.Heartbeat {
	return &tcpkit.Heartbeat{
		Handler:&EchoHeartbeatHandler{},
		Idle: tcpkit.NewHeartbeatIdle(time.Second * 7, time.Second * 15),
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &tcpkit.Config{
		PacketSendChanLimit     : 20,
		PacketReceiveChanLimit  : 20,
		AsyncMessageHand        : true,
	}
	cb := &Callback{}
	protocol := &echo.EchoProtocol{}

	srv := tcpkit.NewServer(cb, protocol)
	srv.SetConfig(config)
	srv.SetHeartbeatMaker(hb_maker)

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
