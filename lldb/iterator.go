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
	//是否包含一下一个
	Next() bool
	Get() string //获取key
	Value() []byte //获取value
	//因为选取范围是左开右闭(]，然而leveldb提供的接口时左闭右开。所以要跳过第一个
	skip()
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

//因为选取范围是左开右闭(]，然而leveldb提供的接口时左闭右开。所以要跳过第一个
func (self *abLLDBIterator) skip() {
	if self.direction == FORWARD {
		if self.it.Next() && self.Get() == self.startKey {

		} else {
			self.it.Prev()
		}
	}
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
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsKV, get:DecodeKV,
	}
	iter.skip()
	return iter
}
//hash key iterator
func NewHashKeyIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsHashSize, get:DecodeHashSize,
	}
	iter.skip()
	return iter
}

//hash label iterator
func NewHashLabelIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsHash, get:DecodeHashLabel,
	}
	iter.skip()
	return iter
}
//queue key iterator
func NewQueueIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsQueueIndex, get:DecodeQueueIndex,
	}
	iter.skip()
	return iter
}
//queue value iterator
func NewQueueValueIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsQueueItem, get:DecodeQueue,
	}
	iter.skip()
	return iter
}
// set list iterator
func NewSetIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is:IsSetSize, get:DecodeSetSize,
	}
	iter.skip()
	return iter
}

//set
type SetIterator struct {
	it    iterator.Iterator
	key   string
	value []byte
}

func (self *SetIterator) skip() {
	if self.it.Next() && self.Get() == self.key {
	} else {
		self.it.Prev()
	}
}

func (self *SetIterator) Next() bool {
	b := self.it.Next() && IsSet(self.it.Key())
	if b {
		key, value := DecodeSet(self.it.Key())
		if self.key != key {
			return false
		}
		self.value = value
	} else {
		_, self.value = "", nil
	}
	return b
}
func (self *SetIterator) Get() string {
	return self.key
}
func (self *SetIterator) Value() []byte {
	return self.value
}
func (self *SetIterator) Release() {
	self.it.Release()
}