package commons

import (
	"errors"
	"fmt"
)

//Try handler(err)
func Try(fun func(), handler ...func(error)) {
	defer Catch(handler...)
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

func SafeExec(fun func()) (err error) {
	Try(fun, func(e error) {
		err = e
	})
	return err
}

func Exec(fn func()) {
	_ = SafeExec(fn)
}

//如果不为空panic错误
func PanicIfPresent(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func Catch(fns ...func(error)) {
	if r := recover(); r != nil && len(fns) > 0 {
		if err, match := r.(error); match {
			for _, fn := range fns {
				fn(err)
			}
		} else {
			err := fmt.Errorf("%v", r)
			for _, fn := range fns {
				fn(err)
			}
		}
	}
}

func CatchError(err error) {
	if r := recover(); r != nil {
		if e, match := r.(error); match {
			err = e
		} else {
			err = fmt.Errorf("%v", r)
		}
	}
}

func Convert(rev interface{}) error {
	if rev == nil {
		return nil
	}
	return fmt.Errorf("%v", rev)
}

//如果不为空，使用msgpanic错误，
func PanicMessageIfPresent(err interface{}, msg string) {
	if err != nil {
		panic(errors.New(fmt.Sprintf(msg, err)))
	}
}
