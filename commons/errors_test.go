package commonKit

import (
    "testing"
)

func TestTry(t *testing.T) {
    i := 1
    j := 1
    Try(func() {
        n := 19 / (i - j)
        t.Log(n)
    }, func(i interface{}) {
        t.Error("已捕获：", i)
    })
}
