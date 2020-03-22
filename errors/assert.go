package errors

import "fmt"

//如果不为空panic错误
func Assert(err error, message ...interface{}) {
	if err != nil {
		if len(message) == 0 {
			panic(Wrap(err, ErrAssert.Error()))
		} else {
			panic(Wrap(err, fmt.Sprint(message...)))
		}
	}
}

func True(checkout bool, message ...interface{}) {
	if !checkout {
		Assert(ErrAssert, message...)
	}
}