package line

import (
	"bufio"
	"fmt"
	"github.com/ihaiker/gokit/remoting"
	"io"
)

type lineCoder struct {
	reader    *bufio.Reader
	LineBreak string //换行分隔符
}

func (l *lineCoder) Encode(channel remoting.Channel, msg interface{}) ([]byte, error) {
	return []byte(fmt.Sprint(msg, l.LineBreak)), nil
}

func (l *lineCoder) Decode(channel remoting.Channel, reader io.Reader) (interface{}, error) {
	if l.reader == nil {
		l.reader = bufio.NewReader(reader)
	}
	ine, _, err := l.reader.ReadLine()
	return string(ine), err
}

func New(lineBreak string) remoting.Coder {
	return &lineCoder{LineBreak: lineBreak}
}
