package lldb

import (
	"github.com/ihaiker/gokit/commons"
)

// key = mark + len(byte key) + byte key + len(value) + value
//
// the key user for key|value = source
//
func EncodeSortedSet(key string, value []byte) []byte {
	out := EncodeSet(key, value)
	out[0] = dt_zset
	return out
}
func DecodeSortedSet(encodeKey []byte) (string, []byte) {
	return DecodeSet(encodeKey)
}

func IsSortedSet(encodeKey []byte) bool {
	return encodeKey[0] == dt_zset
}

func EncodeSortedSetSize(key string) []byte {
	kbytes := []byte(key)
	newBytes := make([]byte, len(kbytes) + 1)
	newBytes[0] = dt_zsize
	copy(newBytes[1:], kbytes)
	return newBytes
}
func DecodeSortedSetSize(encodeKey []byte) string {
	return string(encodeKey[1:])
}

func IsSortedSetSize(encodeKey []byte) bool {
	return encodeKey[0] == dt_zsize
}

func EncodeSortedSetScore(key string, value []byte, score uint64) []byte {
	keyBytes := []byte(key)
	kl := len(keyBytes)
	out := make([]byte, 1 + 1 + kl + 8 /*score byte length*/ + len(value))

	out[0] = dt_zscore           //mark
	out[1] = byte(kl) //key length
	copy(out[2:2 + kl], keyBytes) //key bytes
	commonKit.PutUInt64(out[2 + kl:2 + kl + 8], score)//score
	copy(out[2 + kl + 8:], value)
	return out
}

func DecodeSortedSetScore(encodeByte []byte) (string, []byte, uint64) {
	keyLength := int(encodeByte[1])
	key := make([]byte, keyLength)
	copy(key, encodeByte[2:2 + keyLength])
	score := commonKit.ToUInt64(encodeByte[2 + keyLength:2 + keyLength + 8])
	return string(key), encodeByte[2 + keyLength + 8:], score
}

func IsSortedSetScore(encodeKey []byte) bool {
	return encodeKey[0] == dt_zscore
}