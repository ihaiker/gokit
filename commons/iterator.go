package commons

import "io"

//迭代器
type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type CloseIterator interface {
	Iterator
	io.Closer
}
