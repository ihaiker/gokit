package lldb

/*
	Key/Value编码
*/

//编码kv的key
func EncodeKV(key string) []byte {
	kbytes := []byte(key)
	newBytes := make([]byte, len(kbytes) + 1)
	newBytes[0] = dt_kv
	copy(newBytes[1:], kbytes)
	return newBytes
}
func DecodeKV(encodeKey []byte) string {
	return string(encodeKey[1:])
}

func IsKV(encodeKey [] byte) bool {
	return encodeKey[0] == dt_kv
}
