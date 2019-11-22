package main

import (
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/tlv"
	"github.com/ihaiker/gokit/remoting/example/tlv/msg"
	"github.com/ihaiker/gokit/remoting/handler"
	"log"
	"time"
)

func handlerMaker(ch remoting.Channel) remoting.Handler {
	return handler.Reg().With(handler.Adapter()).
		WithOnMessage(func(session remoting.Channel, message interface{}) {
			logs.Info("接收消息:", message)
			_ = session.Write(msg.NewEcho(fmt.Sprint("ok ", message)))
		})

}

func protocolMaker(ch remoting.Channel) remoting.Coder {
	coder := tlv.NewTLVCoder(1024)
	commons.PanicIfPresent(coder.Reg(msg.NewEcho(time.Now())))
	return coder
}

func main() {
	logs.SetDebugMode(true)
	config := remoting.DefaultTCPConfig()
	config.IdleDuration = 0

	server, err := remoting.NewServer(":6379", config, handlerMaker, protocolMaker)
	if err != nil {
		log.Fatal(err)
	}
	//go func() {
	//	time.Sleep(time.Second * 30)
	//	server.Stop()
	//}()
	server.Start().Wait()
}
