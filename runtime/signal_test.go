package runtimeKit

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"os"
	"testing"
	"time"
)

type NumService struct {
	Num int
}

func (n NumService) Start() error {
	fmt.Println("start num == ", n.Num)
	if n.Num == 3 {
		return os.ErrNotExist
	}
	return nil
}

func (n NumService) Stop() error {
	fmt.Println("stop num == ", n.Num)
	return nil
}

func TestSignal(t *testing.T) {
	logs.SetDebugMode(true)
	lis := NewListener()

	go func() {
		time.Sleep(time.Second)
		lis.Shutdown()
	}()

	lis.Add(&NumService{Num: 1})
	lis.Add(&NumService{Num: 2})
	lis.Add(&NumService{Num: 3})
	lis.AddStop(func() error {
		fmt.Println("4")
		return nil
	})

	t.Log(lis.Wait())
}
