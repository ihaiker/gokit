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