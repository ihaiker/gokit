package errors

import (
	"os"
	"testing"
)

var i, j = 1, 1

func throw() {
	_ = 19 / (i - j)
}

func TestTry(t *testing.T) {
	Try(func() {
		throw()
	}, func(err error) {
		t.Error("已捕获：", err)
	})
}

func TestCatch(t *testing.T) {
	err := func() (err error) {
		defer Catch(func(e error) {
			err = e
		})
		throw()
		return
	}()
	t.Log("运行结果：", err)
}

func C(e error) func(error) {
	return func(err error) {
		e = err
	}
}

func TestSafe(t *testing.T) {
	err := Safe(func() error {
		throw()
		return nil
	})
	t.Log(err)
}
func TestSafeExec(t *testing.T) {
	err := SafeExec(throw)
	t.Log(err)
}

func TestAll(t *testing.T) {
	defer Catch(func(err error) {
		t.Log(err)
	})

	out, err := SafeGet(func() interface{} {
		throw()
		return ""
	})
	True(out == nil)
	True(err != nil)

	out, err = SafeGet(func() interface{} {
		return 1
	})
	True(out == 1)
	Assert(err)

	Assert(os.ErrNotExist, "未发现内容")
}
