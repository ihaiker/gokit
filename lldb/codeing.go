package lldb

import (
	"fmt"
	"encoding/binary"
	"bytes"
)

const (
	dt_kv = 'k'

	//hash
	dt_hash = 'h'; // hashmap(sorted by key)
	dt_hsize = 'H'; // key = size


	//queue
	dt_queue = 'q';
	dt_qsize = 'Q';

	//sset
	dt_sset = 's' //key|vlaue => ""
	dt_ssize = 'S' // key => size

	//zset
	dt_zset = 'x'; // key|value => score
	dt_zscore = 'z'; // key|score => value
	dt_zsize = 'Z'; // key = size
)

// -------------------------- key -------------------------
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

// -------------------------- hash -------------------------
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

// -------------------------- queue -------------------------
func byte2uint64(bs []byte) (uint64, error) {
	buf := bytes.NewBuffer(bs)
	var n uint64
	err := binary.Read(buf, binary.BigEndian, &n)
	return n, err
}

func uint642byte(n uint64) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, n)
	return b_buf.Bytes()
}

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
	return []byte(fmt.Sprintf("%c%020d%s", dt_queue, index, key))
}
func IsQueueItem(encodeKey []byte) bool {
	return encodeKey[0] == dt_queue
}
func DecodeQueue(encodeKey []byte) string {
	return string(encodeKey[1:21])
}

// ------------------------ set ----------------------------
func EncodeSet(key string, value []byte) []byte {
	keyBytes := []byte(key)
	out := make([]byte, 1 + 1 + 1 + len(keyBytes) + len(value))
	out[0] = dt_sset
	out[1] = byte(len(keyBytes))
	out[2] = byte(len(value))
	copy(out[3: 3 + len(keyBytes)], keyBytes)
	copy(out[3 + len(keyBytes):], value)
	return out
}
func DecodeSet(encodeKey []byte) (string, []byte) {
	keyLen := int(encodeKey[1])
	valLen := int(encodeKey[2])
	key := make([]byte, keyLen)
	value := make([]byte, valLen)
	copy(key, encodeKey[3:3 + keyLen])
	copy(value, encodeKey[3 + keyLen:])
	return string(key), value
}
func IsSet(encodeKey []byte) bool {
	return encodeKey[0] == dt_sset
}