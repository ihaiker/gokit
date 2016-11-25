package lldb

import (
	"strings"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Direction int

const (
	FORWARD Direction = 0
	BACKWARD Direction = 1
)

type Iterator interface {
	Next() bool
	Get() string
	Value() []byte
	Release()
}

type abLLDBIterator struct {
	startKey, endKey string
	direction        Direction
	it               iterator.Iterator
	step, limit      int
	is               func([]byte) bool
	get              func([]byte) string
}

func (self *abLLDBIterator) Next() bool {
	if self.step >= self.limit {
		return false
	}
	self.step++

	switch self.direction {
	case FORWARD:
		if self.it.Next() {
			if self.endKey == "" {
				return self.is(self.it.Key())
			} else {
				key := self.Get()
				return strings.Compare(key, self.endKey) <= 0
			}

		}
	case BACKWARD:
		if self.it.Prev() {
			if self.startKey == "" {
				return self.is(self.it.Key())
			} else {
				key := self.Get()
				return strings.Compare(key, self.startKey) > 0
			}
		}
	}
	return false
}
func (self *abLLDBIterator) Release() {
	self.it.Release()
}
func (self *abLLDBIterator) Value() []byte {
	return self.it.Value()
}
func (self *abLLDBIterator) Get() string {
	return self.get(self.it.Key())
}


//key/value iterator
func NewKVIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	return &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsKV, get:DecodeKV,
	}
}
//hash key iterator
func NewHashKeyIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	return &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsHashSize, get:DecodeHashSize,
	}
}

//hash label iterator
func NewHashLabelIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	return &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsHash, get:DecodeHashLabel,
	}
}
//queue key iterator
func NewQueueIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	return &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsQueueIndex, get:DecodeQueueIndex,
	}
}
//queue value iterator
func NewQueueValueIterator(startKey,endKey string, limit int,dir Direction, it iterator.Iterator) Iterator {
	return &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsQueueItem, get:DecodeQueue,
	}
}