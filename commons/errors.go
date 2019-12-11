package commons

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrLimit            = errors.New("limit")         //限制超限
	ErrRange            = errors.New("range overrun") //范围超限
	ErrExists           = errors.New("exists")
	ErrNotSupport       = errors.New("NotSupport") //不支持
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

//安全执行如果出错将被拦截
func Safe(fun func() error) error {
	var err error
	Try(func() {
		err = fun()
	}, func(e error) {
		err = e
	})
	return err
}

func SafeI(fun func() interface{}) interface{} {
	var err interface{}
	Try(func() {
		err = fun()
	}, func(e error) {
		err = e
	})
	return err
}

func SafeExec(fun func()) {
	Try(fun, func(err error) {})
}

func Catch(r interface{}) error {
	var e error = nil
	if r != nil {
		if er, ok := r.(error); ok {
			e = er
		} else if er, ok := r.(string); ok {
			e = errors.New(er)
		} else {
			e = errors.New(fmt.Sprintf("%v", r))
		}
	}
	return e
}

func DCatch(r interface{}, err error) error {
	if err != nil {
		return err
	}
	return Catch(r)
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
		panic(errors.New(fmt.Sprintf(msg, err)))
	}
}
