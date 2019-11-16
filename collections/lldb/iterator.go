package lldb

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"strings"
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

//release the iterator
func (self *abLLDBIterator) Release() {
	self.it.Release()
}
//get the current value
func (self *abLLDBIterator) Value() []byte {
	return self.it.Value()
}
//get the current key
func (self *abLLDBIterator) Get() string {
	return self.get(self.it.Key())
}


//key/value iterator
func newKVIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsKV, get: DecodeKV,
	}
	iter.skip()
	return iter
}
//hash key iterator
func newHashKeyIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsHashSize, get: DecodeHashSize,
	}
	iter.skip()
	return iter
}

//hash label iterator
func newHashLabelIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsHash, get: DecodeHashLabel,
	}
	iter.skip()
	return iter
}
//queue key iterator
func newQueueIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsQueueIndex, get: DecodeQueueIndex,
	}
	iter.skip()
	return iter
}
//queue value iterator
func newQueueValueIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsQueueItem, get: DecodeQueue,
	}
	iter.skip()
	return iter
}
// set set iterator
func newSetIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsSetSize, get: DecodeSetSize,
	}
	iter.skip()
	return iter
}

// set sorted set iterator
func newSortedSetIterator(startKey, endKey string, limit int, dir Direction, it iterator.Iterator) Iterator {
	iter := &abLLDBIterator{
		startKey:startKey, endKey:endKey,
		limit:limit, direction:dir,
		it:it, step:0,
		is: IsSortedSetSize, get: DecodeSortedSetSize,
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

//-- sorted set iterator
type SortedSetIterator struct {
	SetIterator
	score uint64
	step, limit uint64
}
func (self *SortedSetIterator) Score() uint64{
	return self.score
}
func (self *SortedSetIterator) Next() bool {
	if self.step >= self.limit {
		return false
	}
	self.step += 1
	b := self.it.Next() && IsSortedSetScore(self.it.Key())
	if b {
		key, value,score := DecodeSortedSetScore(self.it.Key())
		if self.key != key {
			return false
		}
		self.value = value
		self.score = score
	} else {
		_, self.value = "", nil
	}
	return b
}

func NewSortedSetIterator(it iterator.Iterator, key string, limit uint64) *SortedSetIterator {
	ssi := &SortedSetIterator{}
	ssi.it = it
	ssi.key = key
	ssi.limit = limit
	ssi.step = 0
	return ssi
}
