package future

import (
	"errors"
	"testing"
	"time"
)

func TestDefFuture_Get(t *testing.T) {
	f := Run(func(f Future) (interface{}, error) {
		defer t.Log("step 5")
		t.Log("step 3")
		time.Sleep(time.Second * 6)
		t.Log("step 4")
		return 1, nil
	})
	t.Log("step 2")
	n, err := f.Get()
	t.Log("step 6")
	t.Log(n, err)
}

func TestDefFuture_GetWithTimeout(t *testing.T) {
	t.Log("step 1")
	f := Run(func(f Future) (interface{}, error) {
		defer t.Log("step 5")
		t.Log("step 3")
		time.Sleep(time.Second * 6)
		t.Log("step 4")
		return 1, nil
	})
	t.Log("step 2")
	n, err := f.GetWithTimeout(time.Second * 2)
	t.Log("step 6")
	t.Log(n, err)
	t.Log(f.Get())
	t.Log(f.Get())
	t.Log(f.GetWithTimeout(time.Second * 2))
}

func TestDefFuture_Cancel(t *testing.T) {
	f := Run(func(f Future) (interface{}, error) {
		t.Log("= 1")
		time.Sleep(time.Second * 2)
		t.Log("= 2")
		if f.IsCancelled() {
			return nil, nil
		}
		t.Log("= 3")
		time.Sleep(time.Second * 1)
		t.Log("= 4")
		if f.IsCancelled() {
			return nil, nil
		}
		t.Log("= 5")
		return 6, nil
	})
	t.Log(f.GetWithTimeout(time.Second * 1))
	t.Log(f.Cancel())
	t.Log(f.IsCancelled())
	t.Log(f.Get())
	time.Sleep(time.Second * 3)
}

func TestDefFuture_Exception(t *testing.T) {
	f := Run(func(f Future) (interface{}, error) {
		panic(errors.New("eeeeee"))
	})
	t.Log(f.Get())
}
