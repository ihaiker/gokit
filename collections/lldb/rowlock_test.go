package lldb

import (
	"testing"
)

func TestNewRowLock(t *testing.T) {
	locks,err := NewLocks(2)
	t.Log(locks,err)

	r := locks.Get("123")
	r.RLock()
	t.Log("locking")
	r.RUnlock()
}
