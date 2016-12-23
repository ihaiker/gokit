package lldb

import "github.com/ihaiker/gokit/commons"

func EncodeQueueIndex(key string) []byte {
	kbytes := []byte(key)
	newBytes := make([]byte, len(kbytes) + 1)
	newBytes[0] = dt_qsize
	copy(newBytes[1:], kbytes)
	return newBytes
}
func DecodeQueueIndex(encodeKey []byte) string {
	return string(encodeKey[1:])
}
func IsQueueIndex(encodeKey []byte) bool {
	return encodeKey[0] == dt_qsize
}

//
func EncodeQueue(key string, index uint64) []byte {
	out := make([]byte,1 + 8)
	out[0] = dt_queue
	commonKit.PutUInt64(out[1:],index)
	return out
}
func IsQueueItem(encodeKey []byte) bool {
	return encodeKey[0] == dt_queue
}
func DecodeQueue(encodeKey []byte) string {
	return string(encodeKey[1:])
}

