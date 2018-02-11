package tcpKit

import (
    "bufio"
    "fmt"
    "io"
)

type LineProtocol struct {
    reader    *bufio.Reader
    LineBreak string //换行分隔符
}

func (line *LineProtocol) Encode(msg interface{}) ([]byte, error) {
    return []byte(fmt.Sprintf("%s%s", msg, line.LineBreak)), nil
}

func (line *LineProtocol) Decode(c io.Reader) (interface{}, error) {
    if line.reader == nil {
        line.reader = bufio.NewReader(c)
    }
    ine, _, err := line.reader.ReadLine()
    return string(ine), err
}