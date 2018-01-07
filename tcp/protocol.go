package tcpKit

import (
    "net"
    "fmt"
    "bufio"
)

type Protocol interface {
    Encode(msg interface{}) ([]byte, error)
    Decode(c *net.TCPConn) (interface{}, error)
}

type ProtocolMaker func(c *net.TCPConn) Protocol

type LineProtocol struct {
    reader *bufio.Reader
    Delim  string //换行分隔符
}

func (line *LineProtocol) Encode(msg interface{}) ([]byte, error) {
    return []byte(fmt.Sprintf("%s%s", msg, line.Delim)), nil
}

func (line *LineProtocol) Decode(c *net.TCPConn) (interface{}, error) {
    if line.reader == nil {
        line.reader = bufio.NewReader(c)
    }
    ine, _, err := line.reader.ReadLine()
    return string(ine), err
}
