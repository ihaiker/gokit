package future

import (
    "testing"
    "github.com/ihaiker/gokit/commons/atomic"
    "github.com/magiconair/properties/assert"
    "time"
    "errors"
)

func TestSetFuture_Get(t *testing.T) {
    i := atomic.NewInt()
    f := Set()
    go func() {
        assert.Equal(t, i.IncrementAndGet(), 2, "错误 ")
        f.Set("OVER")
        assert.Equal(t, i.IncrementAndGet(), 3)
    }()
    assert.Equal(t, i.IncrementAndGet(), 1)
    n, err := f.Get()
    assert.Equal(t, i.IncrementAndGet(), 4)
    t.Log(n, err)
}

func TestSetFuture_GetWithTimeout(t *testing.T) {
    i := atomic.NewInt()
    f := Set()
    go func() {
        assert.Equal(t, i.IncrementAndGet(), 2, "错误 ")
        time.Sleep(time.Second * 2)
        if !f.IsCancelled() {
            f.Set("OVER")
        }
        assert.Equal(t, i.IncrementAndGet(), 4)
    }()
    assert.Equal(t, i.IncrementAndGet(), 1)
    t.Log(f.GetWithTimeout(time.Second))
    assert.Equal(t, i.IncrementAndGet(), 3)
    t.Log(f.Cancel())
    t.Log(f.Get())
}


func TestSetFuture_Exception(t *testing.T) {
    i := atomic.NewInt()
    f := Set()
    go func() {
        assert.Equal(t, i.IncrementAndGet(), 2, "错误 ")
        f.Exception(errors.New(" ==== e =========="))
        assert.Equal(t, i.IncrementAndGet(), 3)
    }()
    f.Cancel()
    assert.Equal(t, i.IncrementAndGet(), 1)
    t.Log(f.GetWithTimeout(time.Second))
    assert.Equal(t, i.IncrementAndGet(), 4)
    t.Log(f.Get())
}
