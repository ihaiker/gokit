package main

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/handler"
	"github.com/ihaiker/gokit/remoting/coder/line"
	"time"
)

func main() {
	logs.SetDebugMode(true)

	config := remoting.DefaultTCPConfig()

	clinet, err := remoting.NewClient("127.0.0.1:6379", config, handler.Adapter(), line.New("\n"))
	if err != nil {
		logs.Fatal(err)
		return
	}
	clinet.Start()

	for {
		time.Sleep(time.Millisecond * 10)
		err := clinet.Send(time.Now(), time.Second)
		if err != nil {
			clinet.Close()
			break
		}
	}

	clinet.Wait()
}
