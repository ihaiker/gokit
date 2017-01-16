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


func EncodeQueue(key string, index uint64) []byte {
	keyBytes := []byte(key)
	out := make([]byte, 1 + (1 + len(keyBytes)) + 8)
	out[0] = dt_queue
	out[1] = byte(len(keyBytes))
	copy(out[2:2 + len(keyBytes)], keyBytes)
	commonKit.PutUInt64(out[2 + len(keyBytes):], index)
	return out
}
func IsQueueItem(encodeKey []byte) bool {
	return encodeKey[0] == dt_queue
}
func DecodeQueue(encodeKey []byte) string {
	keyLength := int(encodeKey[1])
	return string(encodeKey[2+keyLength:])
}
func QueueListKey(index uint64) []byte{
	return commonKit.UInt64(index)
}
