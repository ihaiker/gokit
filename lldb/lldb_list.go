package lldb

import (
	"math"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/ihaiker/gokit/commons"
)

const (
	QUEUE_MIN_SEQ uint64 = 10000;
	QUEUE_MAX_SEQ uint64 = math.MaxUint64;
	QUEUE_INIT_SEQ uint64 = QUEUE_MAX_SEQ / 2
)

func (self *LLDBEngine) getQueueIndex(key string) (uint64, uint64, error) {
	val, err := self.data.Get(EncodeQueueIndex(key), self.readOptions)
	if leveldb.ErrNotFound == err {
		return QUEUE_INIT_SEQ, QUEUE_INIT_SEQ + 1, err
	} else {
		min := commonKit.ToUInt64(val[0:8])
		max := commonKit.ToUInt64(val[8:])
		return min, max, nil
	}
}
func (self *LLDBEngine) updateQueueIndex(batch *leveldb.Batch, key string, min, max uint64) {
	if min == QUEUE_INIT_SEQ && max == QUEUE_INIT_SEQ + 1 {
		batch.Delete(EncodeQueueIndex(key))
	} else if min + 1 == max {
		//这个时候队列已经没有字段了
		batch.Delete(EncodeQueueIndex(key))
	} else {
		out := make([]byte, 16)
		commonKit.PutUInt64(out[:8],min)
		commonKit.PutUInt64(out[8:],max)
		batch.Put(EncodeQueueIndex(key), out)
	}
}

func (self *LLDBEngine) _push(dir Direction, key string, value []byte) error {
	minIdx, maxIdx, _ := self.getQueueIndex(key)
	batch := &leveldb.Batch{}
	if BACKWARD == dir {
		batch.Put(EncodeQueue(key, maxIdx), value)
		self.updateQueueIndex(batch, key, minIdx, maxIdx + 1)
	} else {
		batch.Put(EncodeQueue(key, minIdx), value)
		self.updateQueueIndex(batch, key, minIdx - 1, maxIdx)
	}
	return self.data.Write(batch, self.writeOptions)
}

func (self *LLDBEngine) QPush(key string, value []byte) error {
    rwlock := self.queueLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
	return self._push(FORWARD, key, value)
}

func (self *LLDBEngine) QRPush(key string, value []byte) error {
    rwlock := self.queueLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
	return self._push(BACKWARD, key, value)
}

func (self *LLDBEngine) _pop(key string, dir Direction) ([]byte, error) {
	minIdx, maxIdx, err := self.getQueueIndex(key)
	if errors.ErrNotFound == err {
		//not found
		return nil, err
	}
	var index uint64
	if dir == FORWARD {
		index = minIdx + 1
	} else {
		index = maxIdx - 1
	}
	batch := &leveldb.Batch{}
	qKey := EncodeQueue(key, index)
	value, err := self.data.Get(qKey, self.readOptions)
	if err != nil {
		return nil, err
	}
	batch.Delete(qKey)
	if dir == FORWARD {
		self.updateQueueIndex(batch, key, index, maxIdx)
	} else {
		self.updateQueueIndex(batch, key, minIdx, index)
	}
	if err = self.data.Write(batch, self.writeOptions); err != nil {
		return nil, err
	} else {
		return value, err
	}
}

func (self *LLDBEngine) QPop(key string) ([]byte, error) {
    rwlock := self.queueLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
    
	return self._pop(key, FORWARD)
}
func (self *LLDBEngine) QRPop(key string) ([]byte, error) {
    rwlock := self.queueLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
    
	return self._pop(key, BACKWARD)
}

func _queue_index(min, max uint64, index int64) uint64 {
	if index < 0 {
		return max - uint64(0 - index)
	} else {
		return min + 1 + uint64(index)
	}
}

func (self *LLDBEngine) QIndex(key string, index int64) ([]byte, error) {
	minIdx, maxIdx, err := self.getQueueIndex(key)
	if errors.ErrNotFound == err {
		return nil, err
	}
	return self.data.Get(EncodeQueue(key, _queue_index(minIdx, maxIdx, index)), self.readOptions)
}

func (self *LLDBEngine) QList(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeQueueIndex(startKey)
	var endRange []byte = nil
	if (endKey != "" ) {
		endRange = EncodeQueueIndex(endKey + "\001")
	}
	bp := &util.Range{Start:startRange, Limit:endRange}
	it := self.data.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	if it.Valid() {
		it.Seek(startRange) //skip startKey
	} else {
		it.Next()
	}

	return newQueueIterator(startKey, endKey, limit, FORWARD, it), nil
}

func (self *LLDBEngine) QRList(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeQueueIndex(startKey)
	var endRange []byte = nil
	if (endKey != "" ) {
		endRange = EncodeQueueIndex(endKey + "\001")
	} else {
		endRange = EncodeQueueIndex("\255")
	}
	bp := &util.Range{Start:startRange, Limit:endRange}
	it := self.data.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	it.Seek(endRange) //skip startKey
	if ! it.Valid() {
		it.Last()
	} else {
		it.Prev()
	}
	return newQueueIterator(startKey, endKey, limit, BACKWARD, it), nil
}

func (self *LLDBEngine) _trim(key string, limit uint64, dir Direction) (int, error) {
	minIdx, maxIdx, err := self.getQueueIndex(key)
	if errors.ErrNotFound == err {
		//not found
		return 0, err
	}
	deleteNum := uint64(0)
	batch := &leveldb.Batch{}
	if dir == FORWARD {
		for idx := minIdx + 1; idx < maxIdx && idx < minIdx + 1 + limit; idx++ {
			batch.Delete(EncodeQueue(key, idx))
			deleteNum++
		}
		if deleteNum == 0 {
			return 0, nil
		}
		self.updateQueueIndex(batch, key, minIdx + deleteNum, maxIdx)
	} else {
		for idx := maxIdx - 1; idx > minIdx && idx > maxIdx - 1 - limit; idx-- {
			batch.Delete(EncodeQueue(key, idx))
			deleteNum++
		}
		if deleteNum == 0 {
			return 0, nil
		}
		self.updateQueueIndex(batch, key, minIdx, maxIdx - deleteNum)
	}
	return int(deleteNum), self.data.Write(batch, self.writeOptions)
}
//删除队列只剩下limit个
func (self *LLDBEngine) QTrim(key string, limit int) (int, error) {
    rwlock := self.queueLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
    
	return self._trim(key, uint64(limit), FORWARD)
}
//删除对垒只剩下limit个
func (self *LLDBEngine) QRTrim(key string, limit int) (int, error) {
    rwlock := self.queueLock.Get(key)
    rwlock.Lock()
    defer rwlock.Unlock()
    
	return self._trim(key, uint64(limit), BACKWARD)
}

func (self *LLDBEngine) QSize(key string) (uint64, error) {
	minIdx, maxIdx, err := self.getQueueIndex(key)
	//not found
	if errors.ErrNotFound == err {
		return 0, err
	} else {
		return (maxIdx - 1) - (minIdx + 1) + 1, nil
	}
}
func (self *LLDBEngine) QRange(key string, offset, limit uint64) (Iterator, error) {
	if minIdx, maxIdx, err := self.getQueueIndex(key); err != nil {
		return nil, err
	} else {
		start := _queue_index(minIdx, maxIdx, int64(offset))
		end := _queue_index(minIdx, maxIdx, int64(offset + limit))
		if end > maxIdx - 1 {
			end = maxIdx - 1
		}
		startKey := string(QueueListKey(start))
		endKey := string(QueueListKey(end))
		
		startRange := EncodeQueue(key, start)
		//end + 1 是因为levelDB的返回是 [)，详细查阅：util.Range
		endRange := EncodeQueue(key, end + 1)
		bp := &util.Range{Start:startRange, Limit:endRange}
		it := self.data.NewIterator(bp, self.readOptions)
		if err := it.Error(); err != nil {
			return nil, err
		}
		return newQueueValueIterator(startKey, endKey, int(limit), FORWARD, it), nil
	}
}
func (self *LLDBEngine) QSlice(key string, begin, end int64) (Iterator, error) {
	if minIdx, maxIdx, err := self.getQueueIndex(key); err != nil {
		return nil, err
	} else {
		startIdx := _queue_index(minIdx, maxIdx, begin)
		endIdx := _queue_index(minIdx, maxIdx, end)
		if endIdx < startIdx {
			return nil, errors.New("the index start < end")
		}
		startKey := string(QueueListKey(startIdx))
		endKey := string(QueueListKey(endIdx))
		startRange := EncodeQueue(key, startIdx)
		//end + 1 是因为levelDB的返回是 [)，详细查阅：util.Range
		endRange := EncodeQueue(key, endIdx + 1)
		bp := &util.Range{Start:startRange, Limit:endRange}
		it := self.data.NewIterator(bp, self.readOptions)
		if err := it.Error(); err != nil {
			return nil, err
		}
		return newQueueValueIterator(startKey, endKey, int(math.MaxInt8), FORWARD, it), nil
	}
}