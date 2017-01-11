package commonKit

import (
	"encoding/binary"
	"math"
)

var Endian = binary.BigEndian

func UInt8(i uint8) []byte {
	return []byte{byte(i)}
}
func PutUInt8(bs []byte,i uint8) {
	bs[0] = byte(i)
}
func ToUInt8(bs []byte) uint8 {
	return uint8(bs[0])
}
func UInt16(i uint16) []byte {
	bs := make([]byte, 2)
	Endian.PutUint16(bs, i)
	return bs
}
func PutUInt16(bs []byte,i uint16){
	binary.BigEndian.PutUint16(bs,i)
}
func ToUInt16(bs []byte) uint16 {
	return Endian.Uint16(bs)
}
func UInt32(i uint32) []byte {
	bs := make([]byte, 4)
	Endian.PutUint32(bs, i)
	return bs
}
func PutUInt32(bs []byte,i uint32){
	binary.BigEndian.PutUint32(bs,i)
}
func ToUInt32(bs []byte) uint32{
	return Endian.Uint32(bs)
}
func UInt64(i uint64) []byte {
	bs := make([]byte, 8)
	Endian.PutUint64(bs, i)
	return bs
}
func PutUInt64(bs []byte,i uint64){
	binary.BigEndian.PutUint64(bs,i)
}
func ToUInt64(bs []byte) uint64 {
	return Endian.Uint64(bs)
}
func Int8(i int8) []byte {
	return UInt8(uint8(i))
}
func ToInt8(bs []byte) int8 {
	return int8(ToUInt8(bs))
}
func Int16(i uint16) []byte {
	return UInt16(uint16(i))
}
func ToInt16(bs []byte) int16 {
	return int16(ToUInt16(bs))
}
func Int32(i uint32) []byte {
	return UInt32(uint32(i))
}
func ToInt32(bs []byte) int32 {
	return int32(ToUInt32(bs))
}
func Int64(i uint64) []byte {
	return UInt64(uint64(i))
}
func ToInt64(bs []byte) int64 {
	return int64(ToUInt64(bs))
}
func Bool(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}
func ToBool(bs []byte) bool {
	return bs[0] == 1
}
func Float32(f float32) []byte{
	return UInt32(math.Float32bits(f))
}
func ToFloat32(bs []byte) float32 {
	return math.Float32frombits(ToUInt32(bs))
}
func Float64(f float64) []byte{
	return UInt64(math.Float64bits(f))
}
func ToFloat64(bs []byte) float64 {
	return math.Float64frombits(ToUInt64(bs))
}