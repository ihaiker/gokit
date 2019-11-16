package future

import (
	"errors"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetFuture_Get(t *testing.T) {
	i := atomic.NewAtomicInt(0)
	f := Set()
	go func() {
		assert.Equal(t, i.IncrementAndGet(1), 2)
		f.Set("OVER")
		assert.Equal(t, i.IncrementAndGet(1), 3)
	}()
	assert.Equal(t, i.IncrementAndGet(1), 1)
	n, err := f.Get()
	assert.Equal(t, i.IncrementAndGet(1), 4)
	t.Log(n, err)
}

func TestSetFuture_GetWithTimeout(t *testing.T) {
	i := atomic.NewAtomicInt(0)
	f := Set()
	go func() {
		assert.Equal(t, i.IncrementAndGet(1), 2, "错误 ")
		time.Sleep(time.Second * 2)
		if !f.IsCancelled() {
			f.Set("OVER")
		}
		assert.Equal(t, i.IncrementAndGet(1), 4)
	}()
	assert.Equal(t, i.IncrementAndGet(1), 1)
	t.Log(f.GetWithTimeout(time.Second))
	assert.Equal(t, i.IncrementAndGet(1), 3)
	t.Log(f.Cancel())
	t.Log(f.Get())
}

func TestSetFuture_Exception(t *testing.T) {
	i := atomic.NewAtomicInt(0)
	f := Set()
	go func() {
		assert.Equal(t, i.IncrementAndGet(1), 2, "错误 ")
		f.Exception(errors.New(" ==== e =========="))
		assert.Equal(t, i.IncrementAndGet(1), 3)
	}()
	f.Cancel()
	assert.Equal(t, i.IncrementAndGet(1), 1)
	t.Log(f.GetWithTimeout(time.Second))
	assert.Equal(t, i.IncrementAndGet(1), 4)
	t.Log(f.Get())
}
