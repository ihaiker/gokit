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

func TestStack(t *testing.T) {
	defer Catch(func(re error) {
		t.Log(re)
	})
	err := SafeExec(throw)
	Assert(WrapStack(err))
}

func TestIs(t *testing.T) {
	err := Wrap(os.ErrNotExist, "文件不存在")
	t.Log(err)
	is1 := Is(err, os.ErrNotExist)
	t.Log(is1)

	wst := WrapStack(err)
	t.Log(wst)
	is2 := Is(wst, os.ErrNotExist)
	t.Log(is2)

	{
		w3 := Wrap(wst, "在+一层呢")
		t.Log(w3)
		t.Log(Is(w3, os.ErrNotExist))
	}

	{
		w3 := Wrap(err, "在+一层呢")
		t.Log(w3)
		t.Log(Is(w3, os.ErrNotExist))
	}
}
