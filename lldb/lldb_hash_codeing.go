package lldb
/*
    HashMap 编码
*/
// h + byte(len(key)) + []byte(key) + '=' + []byte(label)
func EncodeHash(key, label string) []byte {
	kb, lb := []byte(key), []byte(label)
	kbl, lbl := len(kb), len(lb)
	hkbytes := make([]byte, 1 + 1 + kbl + 1 + lbl)
	hkbytes[0] = dt_hash
	hkbytes[1] = byte(kbl)
	copy(hkbytes[2:], kb)
	hkbytes[1 + 1 + kbl ] = '='
	copy(hkbytes[1 + 1 + kbl + 1 :], lb)
	return hkbytes
}

func DecodeHash(bs []byte) (string, string) {
	length := int(bs[1])
	keyBytes := make([]byte, length)
	copy(keyBytes, bs[2:])
	labelBytes := make([]byte, len(bs) - length - 1 - 1 - 1)
	copy(labelBytes, bs[1 + 1 + length + 1:])
	return string(keyBytes), string(labelBytes)
}
func DecodeHashLabel(bs []byte) string {
	length := int(bs[1])
	labelBytes := make([]byte, len(bs) - length - 1 - 1 - 1)
	copy(labelBytes, bs[1 + 1 + length + 1:])
	return string(labelBytes)
}
func IsHash(bs []byte) bool {
	return bs[0] == dt_hash
}

func EncodeHashSize(key string) []byte {
	kbytes := []byte(key)
	newBytes := make([]byte, len(kbytes) + 1)
	newBytes[0] = dt_hsize
	copy(newBytes[1:], kbytes)
	return newBytes
}
func DecodeHashSize(bs []byte) string {
	return string(bs[1:])
}

func IsHashSize(bs []byte) bool {
	return bs[0] == dt_hsize
}

