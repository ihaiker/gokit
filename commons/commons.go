package commonKit

import "errors"

//Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}
//Try handler(err) and finally
func TryFinally(fun func(), handler func(interface{}), finallyFn func()) {
	defer finallyFn()
	Try(fun,handler)
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