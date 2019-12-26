package executors

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"sync"
	"testing"
	"time"
)

func num(n int, gw *sync.WaitGroup) func() {
	return func() {
		defer gw.Done()
		fmt.Println(time.Now().Format("2006-01-02 15:04:05 .999999"), " num = ", n)
	}
}

func TestNewPool(t *testing.T) {
	logs.SetDebugMode(true)
	gw := new(sync.WaitGroup)
	pool := NewPool(10, 100)

	for i := 0; i < 300; i++ {
		gw.Add(1)
		pool.Add(num(i, gw))
	}
	gw.Wait()
	pool.Shutdown()
}
