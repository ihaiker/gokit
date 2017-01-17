package redis

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
)

type ReplyWriter io.WriterTo

type StatusReply struct {
	code string
}

func (r *StatusReply) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte("+" + r.code + "\r\n"))
	return int64(n), err
}

type IntegerReply struct {
	number int
}

func (r *IntegerReply) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(":" + strconv.Itoa(r.number) + "\r\n"))
	return int64(n), err
}

type BulkReply struct {
	value []byte
}

func writeBytes(value interface{}, w io.Writer) (int64, error) {
	//it's a NullBulkReply
	if value == nil {
		n, err := w.Write([]byte("$-1\r\n"))
		return int64(n), err
	}
	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			n, err := w.Write([]byte("$-1\r\n"))
			return int64(n), err
		}
		wrote, err := w.Write([]byte("$" + strconv.Itoa(len(v)) + "\r\n"))
		if err != nil {
			return int64(wrote), err
		}
		wroteBytes, err := w.Write([]byte(v))
		if err != nil {
			return int64(wrote + wroteBytes), err
		}
		wroteCrLf, err := w.Write([]byte("\r\n"))
		return int64(wrote + wroteBytes + wroteCrLf), err
	case []byte:
		if len(v) == 0 {
			n, err := w.Write([]byte("$-1\r\n"))
			return int64(n), err
		}
		wrote, err := w.Write([]byte("$" + strconv.Itoa(len(v)) + "\r\n"))
		if err != nil {
			return int64(wrote), err
		}
		wroteBytes, err := w.Write(v)
		if err != nil {
			return int64(wrote + wroteBytes), err
		}
		wroteCrLf, err := w.Write([]byte("\r\n"))
		return int64(wrote + wroteBytes + wroteCrLf), err
	case int:
		wrote, err := w.Write([]byte(":" + strconv.Itoa(v) + "\r\n"))
		if err != nil {
			return int64(wrote), err
		}
		return int64(wrote), err
	}

	Debugf("Invalid type sent to writeBytes: %v", reflect.TypeOf(value).Name())
	return 0, errors.New("Invalid type sent to writeBytes")
}

func (r *BulkReply) WriteTo(w io.Writer) (int64, error) {
	return writeBytes(r.value, w)
}

type MonitorReply struct {
	c <-chan string
}

func (r *MonitorReply) WriteTo(w io.Writer) (int64, error) {
	statusReply := &StatusReply{}
	totalBytes := int64(0)
	for line := range r.c {
		statusReply.code = line
		if n, err := statusReply.WriteTo(w); err != nil {
			totalBytes += n
			return int64(totalBytes), err
		} else {
			totalBytes += n
		}
	}
	return totalBytes, nil
}

//for nil reply in multi bulk just set []byte as nil
type MultiBulkReply struct {
	values []interface{}
}

func MultiBulkFromMap(m map[string]interface{}) *MultiBulkReply {
	values := make([]interface{}, len(m)*2)
	i := 0
	for key, val := range m {
		values[i] = []byte(key)
		values[i+1] = val
		i += 2
	}
	return &MultiBulkReply{values: values}
}

func writeMultiBytes(values []interface{}, w io.Writer) (int64, error) {
	if values == nil {
		return 0, errors.New("Nil in multi bulk replies are not ok")
	}
	wrote, err := w.Write([]byte("*" + strconv.Itoa(len(values)) + "\r\n"))
	if err != nil {
		return int64(wrote), err
	}
	wrote64 := int64(wrote)
	for _, v := range values {
		wroteBytes, err := writeBytes(v, w)
		if err != nil {
			return wrote64 + wroteBytes, err
		}
		wrote64 += wroteBytes
	}
	return wrote64, err
}

func (r *MultiBulkReply) WriteTo(w io.Writer) (int64, error) {
	return writeMultiBytes(r.values, w)
}

func ReplyToString(r ReplyWriter) (string, error) {
	var b bytes.Buffer

	_, err := r.WriteTo(&b)
	if err != nil {
		return "ERROR!", err
	}
	return b.String(), nil
}

type FailWriter struct {
	io.ReadWriter
	n int
}

func (fw *FailWriter) Write(buf []byte) (int, error) {
	fw.n -= 1
	if fw.n > 0 {
		return fw.ReadWriter.Write(buf)
	}
	return 0, errors.New("FAILED")
}

// NewFailWriter instanciate a new writer that will fail after n write.
func NewFailWriter(n int) io.ReadWriter {
	w := bytes.NewBuffer([]byte{})
	return &FailWriter{
		ReadWriter: w,
		n:          n,
	}
}