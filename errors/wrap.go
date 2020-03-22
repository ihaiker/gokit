package errors

import "fmt"

type WrapError struct {
	Err     error
	Message string
}

func (w WrapError) Error() string {
	return fmt.Sprintf("%s: %s", w.Message, w.Err)
}

func Wrap(err error, message string) error {
	if _, match := err.(*WrapError); match {
		return err
	} else {
		return &WrapError{
			Err: err, Message: message,
		}
	}
}
func WrapLast(err error, message string) error {
	if we, match := err.(*WrapError); match {
		we.Message = message
		return err
	} else {
		return &WrapError{
			Err: err, Message: message,
		}
	}
}

func Is(err, check error) bool {
	return Root(err) == check
}

func Root(err error) error {
	switch tt := err.(type) {
	case *WrapError:
		return Root(tt.Err)
	case *StackError:
		return Root(tt.Err)
	default:
		return tt
	}
}
