package buffer

import (
	"io"
)

func WriteUInt8(write io.Writer, i uint8) error {
	bs := UInt8(i)
	_, err := write.Write(bs)
	return err
}

func ReadUInt8(reader io.Reader) (uint8, error) {
	bs := make([]byte, 1)
	_, err := io.ReadFull(reader, bs)
	return uint8(bs[0]), err
}

func WriteUInt16(write io.Writer, i uint16) error {
	bs := UInt16(i)
	_, err := write.Write(bs)
	return err
}

func ReadUInt16(reader io.Reader) (uint16, error) {
	bs := make([]byte, 2)
	_, err := io.ReadFull(reader, bs)
	return ToUInt16(bs), err
}

func WriteUInt32(writer io.Writer, i uint32) error {
	bs := UInt32(i)
	_, err := writer.Write(bs)
	return err
}

func ReadUInt32(reader io.Reader) (uint32, error) {
	bs := make([]byte, 4)
	_, err := io.ReadFull(reader, bs)
	return ToUInt32(bs), err
}

func WriteUInt64(writer io.Writer, i uint64) error {
	bs := UInt64(i)
	_, err := writer.Write(bs)
	return err
}

func ReadUInt64(reader io.Reader) (uint64, error) {
	bs := make([]byte, 8)
	_, err := io.ReadFull(reader, bs)
	return ToUInt64(bs), err
}

func WriteInt8(writer io.Writer, i int8) error {
	bs := Int8(i)
	_, err := writer.Write(bs)
	return err
}
func ReadInt8(reader io.Reader) (int8, error) {
	bs := make([]byte, 1)
	_, err := io.ReadFull(reader, bs)
	return ToInt8(bs), err
}

func WriteInt16(write io.Writer, i int16) error {
	bs := Int16(uint16(i))
	_, err := write.Write(bs)
	return err
}

func ReadInt16(reader io.Reader) (int16, error) {
	bs := make([]byte, 2)
	_, err := io.ReadFull(reader, bs)
	return ToInt16(bs), err
}

func WriteInt32(writer io.Writer, i int32) error {
	bs := Int32(uint32(i))
	_, err := writer.Write(bs)
	return err
}

func ReadInt32(reader io.Reader) (int32, error) {
	bs := make([]byte, 4)
	_, err := io.ReadFull(reader, bs)
	return ToInt32(bs), err
}

func WriteInt64(writer io.Writer, i int64) error {
	bs := Int64(uint64(i))
	_, err := writer.Write(bs)
	return err
}

func ReadInt64(reader io.Reader) (int64, error) {
	bs := make([]byte, 8)
	_, err := io.ReadFull(reader, bs)
	return ToInt64(bs), err
}

func WriteBool(write io.Writer, b bool) error {
	if b {
		return WriteInt8(write, 1)
	} else {
		return WriteInt8(write, 0)
	}
}
func ReadBool(reader io.Reader) (bool, error) {
	i, err := ReadInt8(reader)
	return i == 1, err
}

func WriteFloat32(writer io.Writer, f float32) error {
	bs := Float32(f)
	_, err := writer.Write(bs)
	return err
}
func ReadFloat32(reader io.Reader) (float32, error) {
	bs := Float32(0.0)
	_, err := io.ReadFull(reader, bs)
	return ToFloat32(bs), err
}

func WriteFloat64(writer io.Writer, f float64) error {
	bs := Float64(f)
	_, err := writer.Write(bs)
	return err
}
func ReadFloat64(reader io.Reader) (float64, error) {
	bs := Float32(0.0)
	_, err := io.ReadFull(reader, bs)
	return ToFloat64(bs), err
}
