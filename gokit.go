package gokit

import (
	"fmt"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/concurrent/future"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/maths"
	"github.com/ihaiker/gokit/remoting"
	signalKit "github.com/ihaiker/gokit/runtime/signal"
)

//go:generate go run ./concurrent/atomic/genertor/atomic.go ./concurrent/atomic int32 uint32 int64 uint64

func Version() {
	fmt.Println("1.2.0")
	fmt.Println("logs version: ", logs.VERSION)
	fmt.Println("remoting version: ", remoting.VERSION)
	fmt.Println("files version: ", files.VERSION)
	fmt.Println("atomic version:", atomic.VERSION)
	fmt.Println("future version:", future.VERSION)
	fmt.Println("math version:", maths.VERSION)
	fmt.Println("signal version:", signalKit.VERSION)
}
