package buffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

//默认写入的字节流长度，这个长度是因为mark长度标示所限制的。
type ByteSize int

const (
	BYTE8  ByteSize = 8
	BYTE16          = 16
	BYTE32          = 32
	BYTE64          = 64
)

// ---------------------- writer ------------------------------

type ByteWriter struct {
	buffer *bytes.Buffer
	ByteSize
}

func (self *ByteWriter) number(i interface{}) error {
	return binary.Write(self.buffer, binary.BigEndian, i)
}
func (self *ByteWriter) Int8(i int8) error {
	return self.number(i)
}
func (self *ByteWriter) Int16(i int16) error {
	return self.number(i)
}
func (self *ByteWriter) Int32(i int32) error {
	return self.number(i)
}
func (self *ByteWriter) Int64(i int64) error {
	return self.number(i)
}

func (self *ByteWriter) UInt8(i uint8) error {
	return self.number(i)
}
func (self *ByteWriter) UInt16(i uint16) error {
	return self.number(i)
}
func (self *ByteWriter) UInt32(i uint32) error {
	return self.number(i)
}
func (self *ByteWriter) UInt64(i uint64) error {
	return self.number(i)
}
func (self *ByteWriter) Byte(b byte) error {
	return self.UInt8(uint8(b))
}

func (self *ByteWriter) Bool(b bool) error {
	return self.number(b)
}
func (self *ByteWriter) Float32(f float32) error {
	return self.UInt32(math.Float32bits(f))
}
func (self *ByteWriter) Float64(f float64) error {
	return self.UInt64(math.Float64bits(f))
}

func (self *ByteWriter) Write(v []byte) error {
	length := len(v)

	switch self.ByteSize {
	case BYTE8:
		if err := self.UInt8(uint8(length)); err != nil {
			return err
		}
	case BYTE16:
		if err := self.UInt16(uint16(length)); err != nil {
			return err
		}
	case BYTE32:
		if err := self.UInt32(uint32(length)); err != nil {
			return err
		}
	case BYTE64:
		if err := self.UInt64(uint64(length)); err != nil {
			return err
		}
	}

	if length > 0 {
		if _, err := self.buffer.Write(v); err != nil {
			return err
		}
	}
	return nil
}

func (self *ByteWriter) String(v string) error {
	return self.Write([]byte(v))
}

func (self *ByteWriter) ToBytes() []byte {
	return self.buffer.Bytes()
}

func (self *ByteWriter) ToReader() *ByteReader {
	r := NewReader(self.ToBytes())
	r.ByteSize = self.ByteSize
	return r
}

func NewWriter() *ByteWriter {
	return &ByteWriter{
		buffer:   bytes.NewBuffer([]byte{}),
		ByteSize: BYTE16,
	}
}

// ---------------------- reader ------------------------------

type ByteReader struct {
	reader *bytes.Reader
	ByteSize
}

func NewReader(v []byte) *ByteReader {
	return &ByteReader{reader: bytes.NewReader(v), ByteSize: BYTE16}
}

func (self *ByteReader) Int8() (int8, error) {
	return ReadInt8(self.reader)
}
func (self *ByteReader) Int16() (int16, error) {
	return ReadInt16(self.reader)
}
func (self *ByteReader) Int32() (int32, error) {
	return ReadInt32(self.reader)
}
func (self *ByteReader) Int64() (int64, error) {
	return ReadInt64(self.reader)
}

func (self *ByteReader) Byte() (byte, error) {
	i, err := ReadInt8(self.reader)
	return byte(i), err
}

func (self *ByteReader) UInt8() (uint8, error) {
	return ReadUInt8(self.reader)
}
func (self *ByteReader) UInt16() (uint16, error) {
	return ReadUInt16(self.reader)
}
func (self *ByteReader) UInt32() (uint32, error) {
	return ReadUInt32(self.reader)
}
func (self *ByteReader) UInt64() (uint64, error) {
	return ReadUInt64(self.reader)
}

func (self *ByteReader) Float32() (float32, error) {
	return ReadFloat32(self.reader)
}
func (self *ByteReader) Float64() (float64, error) {
	return ReadFloat64(self.reader)
}

func (self *ByteReader) Bytes() ([]byte, error) {
	var len uint64
	var err error
	switch self.ByteSize {
	case BYTE8:
		l, e := self.UInt8()
		len = uint64(l)
		err = e
	case BYTE16:
		l, e := self.UInt16()
		len = uint64(l)
		err = e
	case BYTE32:
		l, e := self.UInt32()
		len = uint64(l)
		err = e
	case BYTE64:
		l, e := self.UInt64()
		len = uint64(l)
		err = e
	}
	if err != nil {
		return nil, err
	}
	if len == 0 {
		return nil, nil
	}
	bs := make([]byte, len)
	rlen, err := io.ReadFull(self.reader, bs)
	if err != nil {
		return nil, err
	}
	if uint64(rlen) != len {
		return nil, errors.New("the length out of index")
	}
	return bs, nil
}

func (self *ByteReader) String() (string, error) {
	bs, err := self.Bytes()
	return string(bs), err
}
