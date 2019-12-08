package rpc

import (
	"errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"strconv"
	"testing"
	"time"
)

func init() {
	logs.SetDebugMode(true)
}

func onMessage(ch remoting.Channel, req *Request) *Response {
	if string(req.Body) == "9" {
		return &Response{id: req.id, Error: errors.New("你错了。你知道吗？这是个世界难题！")}
	} else {
		out := req.URL + " " + string(req.Body)
		return &Response{id: req.id, Body: []byte(out)}
	}
}

func TestRpcServer(t *testing.T) {
	server := NewServer(":6379", onMessage, nil)
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Minute * 10)

	server.Shutdown()
}

func TestRpcClient(t *testing.T) {
	client := NewClient(":6379", onMessage, nil)
	if err := client.Start(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		request := new(Request)
		request.URL = "shutdown"
		request.Body = []byte(strconv.Itoa(i))
		resp := client.Send(request, time.Second*3)
		if resp.Error != nil {
			logger.Debug("error : ", resp.Error)
		} else {
			logger.Debug("response: ", string(resp.Body))
		}
	}
}
