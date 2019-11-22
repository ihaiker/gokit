package asserts

import (
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"reflect"
)

type AssertError struct {
	Code    string `json:"code" yaml:"code" yaml:"code"`
	Message string `json:"message" yaml:"message" yaml:"message"`
}

func (e AssertError) String() string {
	return fmt.Sprint(e.Code, "-: ", e.Message)
}

func newError(code, message string) *AssertError {
	return &AssertError{
		Code:    code,
		Message: message,
	}
}

var Invalid = newError("500", "invalid")

func NotNil(ob interface{}, args ...interface{}) {
	if commons.NotNil(ob) {
		return
	}
	panicIt(args...)
}

func IsTrue(t bool, args ...interface{}) {
	if t == true {
		return
	}
	panicIt(args...)
}

func panicIt(args ...interface{}) {
	switch len(args) {
	case 0:
		panic(Invalid)
	case 1:
		switch args[0].(type) {
		case string:
			panic(newError("500", args[0].(string)))
		case *AssertError, AssertError:
			panic(args[0])
		case error:
			panic(newError("500", args[0].(error).Error()))
		}
	case 2:
		if reflect.TypeOf(args[0]).String() == "string" && reflect.TypeOf(args[1]).String() == "string" {
			panic(newError(args[0].(string), args[1].(string)))
		}
	}
	panic(errors.New(fmt.Sprint(args...)))
}
