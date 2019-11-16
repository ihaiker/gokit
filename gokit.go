package gokit

import (
	"fmt"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/concurrent/future"
	fileKit "github.com/ihaiker/gokit/files"
	logsKit "github.com/ihaiker/gokit/logs"
	tcpKit "github.com/ihaiker/gokit/tcp"
)

//go:generate go run ./concurrent/atomic/genertor/atomic.go ./concurrent/atomic int32 uint32 int64 uint64

func Version() {
	fmt.Println("1.0.0")
	fmt.Println("logs version: ", logsKit.VERSION)
	fmt.Println("tcpKit version: ", tcpKit.VERSION)
	fmt.Println("fileKit version: ", fileKit.VERSION)
	fmt.Println("atomicKit version:", atomic.VERSION)
	fmt.Println("future version:", future.VERSION)
}
