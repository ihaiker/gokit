package atomic

import "sync/atomic"

type AtomicInt64 struct {
	value int64
}

func (self *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&self.value)
}

func (self *AtomicInt64) IncrementAndGet() int64 {
	return atomic.AddInt64(&self.value, 1)
}

func (self *AtomicInt64) GetAndIncrement() (int64) {
	return self.GetAndAdd(1)
}

func (self *AtomicInt64) DecrementAndGet() (int64) {
	return self.AddAndGet(-1)
}

func (self *AtomicInt64) GetAndDecrement() (int64) {
	return self.GetAndAdd(-1)
}

func (self *AtomicInt64) AddAndGet(i int64) (int64) {
	return atomic.AddInt64(&self.value, i)
}

func (self *AtomicInt64) GetAndAdd(i int64) (int64) {
	var ret int64
	for ; ; {
		ret = atomic.LoadInt64(&self.value)
		if atomic.CompareAndSwapInt64(&self.value, ret, ret + i) {
			break
		}
	}
	return ret
}

func (self *AtomicInt64) CompareAndSet(expect int64, update  int64) (bool) {
	return atomic.CompareAndSwapInt64(&self.value, expect, update)
}