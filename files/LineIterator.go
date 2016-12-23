package fileKit

import (
	"github.com/ihaiker/gokit/commons"
	"bufio"
	"os"
)

type LineIterator struct {
	commonKit.CloseIterator
	reader  *bufio.Reader
	file    *os.File
	current []byte
}

func (self *LineIterator) HasNext() bool {
	line, _, err := self.reader.ReadLine()
	if err != nil {
		self.current = nil
		return false
	}
	self.current = line
	return true
}

func (self *LineIterator) Next() interface{} {
	defer func() {
		self.current = nil
	}()
	return self.current
}

func (self *LineIterator) Close() error {
	return self.file.Close()
}

func newIterator(file *os.File) commonKit.CloseIterator {
	return &LineIterator{
		file:file,
		reader:bufio.NewReader(file),
	}
}

