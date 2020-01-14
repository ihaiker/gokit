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
	return handler.Reg().With(handler.Adapter())
}

func protocolMaker(ch remoting.Channel) remoting.Coder {
	return line.New("\n")
}

func main() {
	logs.SetDebugMode(true)

	config := remoting.DefaultOptions()
	server := remoting.NewServer(":6379", config, handlerMaker, protocolMaker)

	err := server.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	runtimeKit.WaitTC(time.Second, func() {
		_ = server.Stop()
		server.Wait()
	})
}
