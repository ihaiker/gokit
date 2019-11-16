package fileKit

import (
	"github.com/ihaiker/gokit/commons"
	"bufio"
	"os"
)

type lineIterator struct {
	commons.CloseIterator
	reader  *bufio.Reader
	file    *os.File
	current []byte
}

func (self *lineIterator) HasNext() bool {
	line, _, err := self.reader.ReadLine()
	if err != nil {
		self.current = nil
		return false
	}
	self.current = line
	return true
}

func (self *lineIterator) Next() interface{} {
	defer func() {
		self.current = nil
	}()
	return self.current
}

func (self *lineIterator) Close() error {
	return self.file.Close()
}

func newIterator(file *os.File) commons.CloseIterator {
	return &lineIterator{
		file:   file,
		reader: bufio.NewReader(file),
	}
}
