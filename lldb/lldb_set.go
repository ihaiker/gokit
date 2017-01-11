package lldb

import (
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/rand"
	"time"
	"github.com/ihaiker/gokit/commons"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func _add_set(self *LLDBEngine, batch *leveldb.Batch, key string, value []byte) (int, error) {
	encodeKey := EncodeSet(key, value)
	_, err := self.data.Get(encodeKey, self.readOptions)
	if err != nil {
		if err == errors.ErrNotFound {
			batch.Put(encodeKey, []byte{})
			return 1, nil
		} else {
			return 0, err
		}
	}
	return 0, nil
}

func _del_set(self *LLDBEngine, batch *leveldb.Batch, key string, value []byte) (int, error) {
	encodeKey := EncodeSet(key, value)
	_, err := self.data.Get(encodeKey, self.readOptions)
	if err != nil {
		if err == errors.ErrNotFound {
			return 0, nil
		} else {
			return 0, err
		}
	}
	batch.Delete(encodeKey)
	return 1, nil
}

func _update_sset_size(self *LLDBEngine, batch *leveldb.Batch, key string, incr int64) (uint64, error) {
	size, err := self.SSize(key)
	if err != nil {
		return 0, err
	}
	if incr < 0 {
		size = size - uint64(0 - incr)
	} else {
		size = size + uint64(incr)
	}
	if size <= 0 {
		batch.Delete(EncodeSetSize(key))
		return 0, nil
	} else {
		batch.Put(EncodeSetSize(key), commonKit.UInt64(size))
		return size, nil
	}
}

func (self *LLDBEngine) SAdd(key string, value []byte) (int, error) {
	rwlock := self.setLock.Get(key)
	rwlock.Lock()
	defer rwlock.Unlock()
	
	batch := &leveldb.Batch{}
	insert, err := _add_set(self, batch, key, value)
	if err != nil {
		return 0, err
	} else if insert == 0 {
		return 0, nil
	}
	if _, err = _update_sset_size(self, batch, key, 1); err != nil {
		return 0, err
	}
	return insert, self.data.Write(batch, self.writeOptions)
}

func (self *LLDBEngine) SDel(key string, value []byte) (int, error) {
	rwlock := self.setLock.Get(key)
	rwlock.Lock()
	defer rwlock.Unlock()
	
	batch := &leveldb.Batch{}
	del, err := _del_set(self, batch, key, value)
	if err != nil || del == 0 {
		return del, err
	}
	if _, err := _update_sset_size(self, batch, key, -1); err != nil {
		return 0, err
	}

	return 1, self.data.Write(batch, self.writeOptions)

}

func (self *LLDBEngine) SExits(key string, value[]byte) (bool, error) {
	_, err := self.data.Get(EncodeSet(key, value), self.readOptions)
	if err == errors.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func (self *LLDBEngine) SMembers(key string) (Iterator, error) {
	startRange := EncodeSet(key, []byte{})
	it := self.data.NewIterator(&util.Range{Start:startRange}, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	return &SetIterator{it:it, key:key}, nil
}

func (self *LLDBEngine) SRandomMember(key string) ([]byte, error) {
	size, err := self.SSize(key)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return nil, nil
	}
	idx := rand.Uint64() % size

	it := self.data.NewIterator(&util.Range{Start:EncodeSet(key, []byte{})}, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	var i uint64 = 0
	for ; i < idx && it.Next(); i++ {
		//k := it.Key()
		//v := it.Value()
	}
	if it.Next() {
		_, v := DecodeSet(it.Key())
		return v, nil
	}
	return nil, nil
}

func (self *LLDBEngine) SPop(key string) ([]byte, error) {
	rwlock := self.setLock.Get(key)
	rwlock.Lock()
	defer rwlock.Unlock()
	
	val, err := self.SRandomMember(key)
	if err != nil {
		return nil, err
	}
	_, err = self.SDel(key, val)
	return val, err
}

func (self *LLDBEngine) SList(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeSetSize(startKey)
	var endRange []byte = nil
	if (endKey != "" ) {
		endRange = EncodeSetSize(endKey + "\001")
	}
	it := self.data.NewIterator(&util.Range{Start:startRange, Limit:endRange}, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	return newSetIterator(startKey, endKey, limit, FORWARD, it), nil
}

func (self *LLDBEngine) SSize(key string) (uint64, error) {
	val, err := self.data.Get(EncodeSetSize(key), self.readOptions)
	if err != nil {
		if err == errors.ErrNotFound {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		return commonKit.ToUInt64(val), nil
	}
}