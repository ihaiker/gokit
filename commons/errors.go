package commonKit

import (
    "errors"
    "fmt"
)

//Try handler(err)
func Try(fun func(), handler func(error)) {
    defer func() {
        if err := Catch(recover()); err != nil {
            handler(err)
        }
    }()
    fun()
}

//Try handler(err) and finally
func TryFinally(fun func(), handler func(error), finallyFn func()) {
    defer finallyFn()
    Try(fun, handler)
}

func Catch(r interface{}) error {
    var e error = nil
    if r != nil {
        if er, ok := r.(error); ok {
            e = er
        } else if er, ok := r.(string); ok {
            e = errors.New(er)
        } else {
            e = errors.New(fmt.Sprintf("%s", r))
        }
    }
    return e
}

//如果不为空panic错误
func IfPanic(err interface{}) {
    if err != nil {
        panic(err)
    }
}

//如果不为空panic错误
func PanicIfPresent(err interface{}) {
    if err != nil {
        panic(err)
    }
}

//如果不为空，使用msgpanic错误，
func PanicMessageIfPresent(err interface{}, msg string) {
    if err != nil {
        panic(errors.New(msg))
    }
}
