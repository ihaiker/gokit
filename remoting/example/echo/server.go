package main

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/line"
	"github.com/ihaiker/gokit/remoting/handler"
	signalKit "github.com/ihaiker/gokit/runtime/signal"
	"log"
	"os"
	"time"
)

func handlerMaker(ch remoting.Channel) remoting.Handler {
	return handler.Reg().With(handler.Adapter()).
		WithOnMessage(func(session remoting.Channel, msg interface{}) {
			logs.Debug("新消息：", msg)
			_ = session.Write(fmt.Sprint("you see: ", msg))
			session.AsyncWrite(fmt.Sprint("you see: ", msg, ", are you sure?"), time.Second, func(msg interface{}, err error) {
				if err != nil {
					logs.Error("异步消息异常:", msg, " error:", err)
				} else {
					logs.Debug("异步消息发送成功！", msg)
				}
			})
		}).
		WithOnIdle(func(session remoting.Channel) {
			_ = session.Write("PING")
		})

}

func protocolMaker(ch remoting.Channel) remoting.Coder {
	return line.New("\n")
}

func main() {
	logs.SetDebugMode(true)

	config := remoting.DefaultTCPConfig()
	server, err := remoting.NewServer("unix://tmp/test.sock", config, handlerMaker, protocolMaker)
	if err != nil {
		log.Fatal(err)
	}
	//go func() {
	//	time.Sleep(time.Second * 30)
	//	server.Stop()
	//}()

	server.Start()

	signalKit.Signal(nil, func(signal ...os.Signal) {
		server.Stop().Wait()
	})
}
