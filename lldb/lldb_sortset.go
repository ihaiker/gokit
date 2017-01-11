package lldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/ihaiker/gokit/commons"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"math"
	"bytes"
)

func (self *LLDBEngine) _zset_one(batch *leveldb.Batch, key string, value []byte, score uint64) (int, error) {
	oldSource, err := self.ZScore(key, value)
	//查询出错
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}
	found := !(err == leveldb.ErrNotFound)
	if found {
		if oldSource == score {
			return 0, nil
		}
		batch.Delete(EncodeSortedSetScore(key, value, oldSource))
	}
	batch.Put(EncodeSortedSet(key, value), commonKit.UInt64(score))
	batch.Put(EncodeSortedSetScore(key, value, score), []byte{})
	return commonKit.IfElse(found, 0, 1).(int), nil
}

func (self *LLDBEngine) _incr_sorted_size(batch *leveldb.Batch, key string, incr int64) error {
    rwlock := self.sortedSetLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
    
	if size, err := self.ZSize(key); err != nil {
		return err
	} else {
		if incr < 0 {
			size = size - uint64(0 - incr)
		} else {
			size = size + uint64(incr)
		}
		if size <= 0 {
			batch.Delete(EncodeSortedSetSize(key))
		} else {
			batch.Put(EncodeSortedSetSize(key), commonKit.UInt64(size))
		}
		return nil
	}
}

//添加一项有序
func (self *LLDBEngine) ZAdd(key string, value []byte, score uint64) (int, error) {
	batch := &leveldb.Batch{}
	ret, err := self._zset_one(batch, key, value, score);
	if err != nil {
		return 0, err
	} else {
		if ret > 0 {
			if err := self._incr_sorted_size(batch, key, int64(1)); err != nil {
				return 0, err
			}
		}
		return ret, self.data.Write(batch, self.writeOptions)
	}
}
//获取有序set的权重
func (self *LLDBEngine) ZScore(key string, value []byte) (uint64, error) {
	if v, err := self.data.Get(EncodeSortedSet(key, value), self.readOptions); err != nil {
		return 0, err
	} else {
		return commonKit.ToUInt64(v), nil
	}
}

//获取key的当前元素数
func (self *LLDBEngine) ZSize(key string) (uint64, error) {
	v, err := self.data.Get(EncodeSortedSetSize(key), self.readOptions)
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	} else if err == leveldb.ErrNotFound {
		return 0, nil
	} else {
		return commonKit.ToUInt64(v), nil
	}
}

func (self *LLDBEngine) ZIncrBy(key string, value []byte, score uint64) (uint64, error) {
	oldSource, err := self.ZScore(key, value)
	//查询出错
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}
	_, err = self.ZAdd(key, value, oldSource + score)
	return oldSource + score, err
}

//redis zrem
func (self *LLDBEngine) ZDel(key string, value []byte) (int, error) {
	score, err := self.ZScore(key, value)
	//查询出错
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}
	//没有当前属性
	if err == leveldb.ErrNotFound {
		return 0, nil
	}
	batch := &leveldb.Batch{}
	batch.Delete(EncodeSortedSet(key, value))
	batch.Delete(EncodeSortedSetScore(key, value, score))
	if err := self._incr_sorted_size(batch, key, -1); err != nil {
		return 0, err
	}
	return 1, self.data.Write(batch, self.writeOptions)
}

func (self *LLDBEngine) iterator(key string, minScore, maxScore uint64) (iterator.Iterator, error) {
	it := self.data.NewIterator(&util.Range{
		Start: EncodeSortedSetScore(key, []byte{}, minScore),
		Limit:EncodeSortedSetScore(key, []byte{255}, maxScore),
	}, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	return it, nil
}

//返回处于区间 [start,end] key 数量.
func (self *LLDBEngine) ZCount(key string, min, max uint64) (uint64, error) {
	it, err := self.iterator(key, min, max)
	if err != nil {
		return 0, err
	}
	defer it.Release()

	count := uint64(0)
	for ; it.Next(); {
		count++
	}
	return count, nil
}
//列出名字处于区间 (name_start, name_end] 的 zset.
func (self *LLDBEngine) ZList(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeSortedSetSize(startKey)
	var endRange []byte = nil
	if (endKey != "" ) {
		endRange = EncodeSortedSetSize(endKey + "\255")
	}
	it := self.data.NewIterator(&util.Range{Start:startRange, Limit:endRange}, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	return newSortedSetIterator(startKey, endKey, limit, FORWARD, it), nil
}

//清空key
func (self *LLDBEngine) ZClear(key string) (uint64, error) {
	if it, err := self.ZRange(key, 0, math.MaxUint64, math.MaxUint64); err == nil {
		batch := &leveldb.Batch{}
		del := uint64(0)
		for ; it.Next(); {
			value := it.Value()
			score := it.Score()
			batch.Delete(EncodeSortedSet(key, value))
			batch.Delete(EncodeSortedSetScore(key, value, score))
			del += 1
		}
		batch.Delete(EncodeSortedSetSize(key))
		return del, self.data.Write(batch, self.writeOptions)
	} else {
		return 0, err
	}
}
func (self *LLDBEngine) ZRange(key string, minScore, maxScore, limit uint64) (*SortedSetIterator, error) {
	it, err := self.iterator(key, minScore, maxScore)
	if err != nil {
		return nil, err
	}
	return NewSortedSetIterator(it, key, limit), nil
}
func (self *LLDBEngine) ZRank(key string, value []byte) (uint64, error) {
	it, err := self.ZRange(key, 0, math.MaxUint64, math.MaxUint64)
	if err != nil {
		return 0, err
	}
	defer it.Release()
	idx := uint64(0)
	for ; it.Next(); {
		if bytes.Equal(it.Value(), value) {
			return idx, nil
		}
		idx += 1
	}
	return idx, nil
}

func (self *LLDBEngine) ZDelByRank(key string, startRank, stopRank uint64) (uint64, error) {
	if it, err := self.ZRange(key, 0, math.MaxUint64, math.MaxUint64); err != nil {
		return 0, err
	} else {
		del := uint64(0)
		batch := &leveldb.Batch{}
		for ; it.Next() && (del >= startRank && del <= stopRank); {
			value := it.Value()
			score := it.Score()
			batch.Delete(EncodeSortedSet(key, value))
			batch.Delete(EncodeSortedSetScore(key, value, score))
			del += 1
		}
		if del != 0 {
			self._incr_sorted_size(batch, key, int64(0) - int64(del))
		}
		return del, self.data.Write(batch, self.writeOptions)
	}
}

func (self *LLDBEngine) ZDelByScore(key string, startScore, stopScore uint64) (uint64, error) {
	if it, err := self.ZRange(key, startScore, stopScore, math.MaxUint64); err != nil {
		return 0, err
	} else {
		del := uint64(0)
		batch := &leveldb.Batch{}
		for ; it.Next(); {
			value := it.Value()
			score := it.Score()
			batch.Delete(EncodeSortedSet(key, value))
			batch.Delete(EncodeSortedSetScore(key, value, score))
			del += 1
		}
		if del != 0 {
			self._incr_sorted_size(batch, key, int64(0) - int64(del))
		}
		return del, self.data.Write(batch, self.writeOptions)
	}
}