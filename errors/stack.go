package errors

import "fmt"

type StackError struct {
	Err   error
	Stack string
}

func (w StackError) Error() string {
	return fmt.Sprintf("%s \n%s", w.Err, w.Stack)
}

func WrapStack(err error) error {
	if stack, match := err.(*StackError); match {
		stack.Stack += "\n " + err.Error()
		return stack
	} else {
		return &StackError{
			Err: err, Stack: Stack(),
		}
	}
}
