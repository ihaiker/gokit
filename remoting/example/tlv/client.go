package main

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/tlv"
	"github.com/ihaiker/gokit/remoting/example/tlv/msg"
	"github.com/ihaiker/gokit/remoting/handler"
	"math/rand"
	"time"
)

func main() {
	defer logs.CloseAll()
	logs.SetDebugMode(true)

	reg := tlv.NewTLVCoder(1024)
	_ = reg.Reg(msg.NewEcho(time.Now()))

	config := remoting.DefaultOptions()
	config.IdleTimeSeconds = 0

	clinet := remoting.NewClient("127.0.0.1:6379", config, handler.Adapter(), reg)

	err := clinet.Start()
	if err != nil {
		logs.Fatal(err)
		return
	}

	for {
		n := rand.Intn(1000)
		time.Sleep(time.Duration(n) * time.Millisecond)
		err := clinet.Write(msg.NewEcho(n), time.Second)
		if err != nil {
			logs.Error("消息错误：", err)
			clinet.Close()
			break
		}
	}

	clinet.Wait()
}
