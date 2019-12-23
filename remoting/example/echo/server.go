package main

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/line"
	"github.com/ihaiker/gokit/remoting/handler"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"os"
	"time"
)

func handlerMaker(ch remoting.Channel) remoting.Handler {
	return handler.Reg().With(handler.Adapter()).
		WithOnMessage(func(session remoting.Channel, msg interface{}) {
			logs.Debug("新消息：", msg)
			_ = session.Write(fmt.Sprint("you see: ", msg), time.Second)
			session.AsyncWrite(fmt.Sprint("you see: ", msg, ", are you sure?"), time.Second, func(msg interface{}, err error) {
				if err != nil {
					logs.Error("异步消息异常:", msg, " error:", err)
				} else {
					logs.Debug("异步消息发送成功！", msg)
				}
			})
		}).
		WithOnIdle(func(session remoting.Channel) {
			_ = session.Write("PING", time.Second)
		})

}

func protocolMaker(ch remoting.Channel) remoting.Coder {
	return line.New("\n")
}

func main() {
	logs.SetDebugMode(true)

	config := remoting.DefaultTCPConfig()
	server := remoting.NewServer(":6379", config, handlerMaker, protocolMaker)

	err := server.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runtimeKit.NewListener().WaitWithTimeout(time.Second, func() {
		_ = server.Stop()
		server.Wait()
	})
}
