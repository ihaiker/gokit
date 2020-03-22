package main

import (
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/tlv"
	"github.com/ihaiker/gokit/remoting/example/tlv/msg"
	"github.com/ihaiker/gokit/remoting/handler"
	"log"
	"reflect"
	"time"
)

func handlerMaker(ch remoting.Channel) remoting.Handler {
	return handler.Reg().With(handler.Adapter()).On(remoting.MessageEvent, func(event *remoting.Event) {
		num := event.Values[0]
		fmt.Println("收到消息:", num, ",", reflect.TypeOf(num).String())
	})
}

func protocolMaker(ch remoting.Channel) remoting.Coder {
	coder := tlv.NewTLVCoder(1024)
	errors.Assert(coder.Reg(msg.NewEcho(time.Now())))
	return coder
}

func main() {
	logs.SetDebugMode(true)
	config := remoting.DefaultOptions()
	config.IdleTimeSeconds = 0

	server := remoting.NewServer(":6379", config, handlerMaker, protocolMaker)

	//go func() {
	//	time.Sleep(time.Second * 30)
	//	server.Stop()
	//}()

	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
	server.Wait()
}
