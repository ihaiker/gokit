package commons

import (
    "errors"
    "fmt"
    "testing"
)

func TestTry(t *testing.T) {
    i := 1
    j := 1
    Try(func() {
        n := 19 / (i - j)
        t.Log(n)
    }, func(err error) {
        t.Error("已捕获：", err)
    })
}

func TestCatch(t *testing.T) {
    testCatch := func() (err error) {
        defer Catch(func(e error) {
            err = e
        })

        i := 1
        j := 1
        n := 19 / (i - j)
        fmt.Println(n)
        return
    }
    err := testCatch()
    t.Log(err)
}

func TestDCatch(t *testing.T) {
    testCatch := func() (err error) {
        defer Catch(func(e error) {
            err = e
        })

        err = errors.New("test")

        i := 1
        j := 1
        n := 19 / (i - j)
        fmt.Println(n)
        return
    }
    err := testCatch()
    t.Log(err)
}
