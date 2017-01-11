package lldb


// ------------------------ set ----------------------------
func EncodeSet(key string, value []byte) []byte {
	keyBytes := []byte(key)
	out := make([]byte, 1 + 1 + len(keyBytes) + len(value))

	out[0] = dt_sset
	out[1] = byte(len(keyBytes))

	copy(out[2: 2 + len(keyBytes)], keyBytes)
	copy(out[2 + len(keyBytes):], value)
	return out
}
func DecodeSet(encodeKey []byte) (string, []byte) {
	keyLen := int(encodeKey[1])
	key := make([]byte, keyLen)
	value := make([]byte, len(encodeKey) - 2 - keyLen)
	copy(key, encodeKey[2:2 + keyLen])
	copy(value, encodeKey[2 + keyLen:])
	return string(key), value
}
func IsSet(encodeKey []byte) bool {
	return encodeKey[0] == dt_sset
}

func EncodeSetSize(key string) []byte {
	keyByes := []byte(key)
	bs := make([]byte, len(keyByes) + 1)
	bs[0] = dt_ssize
	copy(bs[1:], keyByes)
	return bs
}
func IsSetSize(encodeKey []byte) bool {
	return encodeKey[0] == dt_ssize
}
func DecodeSetSize(encodeKey []byte) string {
	return string(encodeKey[1:])
}
