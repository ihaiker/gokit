package main

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/line"
	"github.com/ihaiker/gokit/remoting/handler"
	"time"
)

func main() {
	logs.SetDebugMode(true)

	config := remoting.DefaultOptions()

	clinet := remoting.NewClient("127.0.0.1:6379", config, handler.Adapter(), line.New("\n"))

	if err := clinet.Start(); err != nil {
		logs.Fatal(err)
		return
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		err := clinet.Write(time.Now(), time.Second)
		if err != nil {
			break
		}
	}

	_ = clinet.Close()
	clinet.Wait()
	fmt.Println("OVER")
}
