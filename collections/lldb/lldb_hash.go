package lldb

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math"
	"strconv"
)

func (self *LLDBEngine) incr_hsize(batch *leveldb.Batch, key string, incr int) {
	size := self.HSize(key)
	size = size + incr
	if size == 0 {
		batch.Delete(EncodeHashSize(key))
	} else {
		batch.Put(EncodeHashSize(key), []byte(strconv.Itoa(size)))
	}
}

func (self *LLDBEngine) hset_one(batch *leveldb.Batch, key, label string, value []byte) int {
	val, err := self.HGet(key, label)
	if leveldb.ErrNotFound == err {
		batch.Put(EncodeHash(key, label), value)
		return 1
	} else if !bytes.Equal(val, value) {
		batch.Put(EncodeHash(key, label), value)
		return 1
	} else {
		return 0
	}
}

func (self *LLDBEngine) hdel_one(batch *leveldb.Batch, key, label string) int {
	_, err := self.HGet(key, label)
	if leveldb.ErrNotFound == err {
		return 0
	} else {
		batch.Delete(EncodeHash(key, label))
		return 1
	}
}

//set key value
//设置指定 key 的值内容.
func (self *LLDBEngine) HSet(key, label string, value []byte) (int, error) {
	rwlock := self.setLock.Get(key)
	rwlock.Lock()
	defer rwlock.Unlock()

	batch := &leveldb.Batch{}
	insert := self.hset_one(batch, key, label, value)
	if insert == 0 {
		return 0, nil
	}
	self.incr_hsize(batch, key, 1)
	return insert, self.data.Write(batch, self.writeOptions)
}

func (self *LLDBEngine) HGet(key, label string) ([]byte, error) {
	return self.data.Get(EncodeHash(key, label), self.readOptions)
}
func (self *LLDBEngine) HDel(key, label string) (int, error) {
	rwlock := self.setLock.Get(key)
	rwlock.Lock()
	defer rwlock.Unlock()

	batch := &leveldb.Batch{}
	del := self.hdel_one(batch, key, label)
	if del == 0 {
		return 0, nil
	}
	self.incr_hsize(batch, key, -1)
	return del, self.data.Write(batch, self.writeOptions)
}

//get the hash size
func (self *LLDBEngine) HSize(key string) int {
	val, err := self.data.Get(EncodeHashSize(key), self.readOptions)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0
		} else {
			return -1
		}
	} else {
		if val == nil {
			return -1
		} else {
			o, err := strconv.ParseInt(string(val), 10, 8)
			if err != nil {
				return -1
			} else {
				return int(o)
			}
		}
	}
	return 0
}

//list the name of `hash` data
func (self *LLDBEngine) HList(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeHashSize(startKey)
	var endRange []byte = nil
	if (endKey != "" ) {
		endRange = EncodeHashSize(endKey + "\001")
	}
	bp := &util.Range{Start:startRange, Limit:endRange}

	it := self.data.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	it.Seek(startRange) //skip startKey
	if (!it.Valid()) {
		it.Last();
	}
	return newHashKeyIterator(startKey, endKey, limit, FORWARD, it), nil
}
//list from left to right the name of `hash` data
func (self *LLDBEngine) HRList(startKey, endKey string, limit int) (Iterator, error) {
	var startRange, endRange []byte
	startRange = EncodeHashSize(startKey)
	if endKey != "" {
		endRange = EncodeHashSize(endKey + "\001")
	} else {
		endRange = EncodeHashSize("\255")
	}
	bp := &util.Range{Start:startRange, Limit:endRange}

	it := self.data.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	it.Last()
	it.Seek(endRange)

	return newHashKeyIterator(startKey, endKey, limit, BACKWARD, it), nil
}

//scan the hash table data by the give key
func (self *LLDBEngine) HScan(key string, startLabel, endLabel string, limit int) (Iterator, error) {
	startRange := EncodeHash(key, startLabel)
	var endRange []byte = nil
	if endLabel != "" {
		endRange = EncodeHash(key, endLabel + "\001")
	} else {
		endRange = EncodeHash(key, "\255")
	}

	bp := &util.Range{Start:startRange, Limit:endRange}
	it := self.data.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	it.Seek(startRange) //skip startKey
	if (!it.Valid()) {
		it.Last();
	}
	return newHashLabelIterator(startLabel, endLabel, limit, FORWARD, it), nil
}

func (self *LLDBEngine) HGetAll(key string) (Iterator, error) {
	return self.HScan(key, "", "", math.MaxInt32)
}
