package atomic

import "sync/atomic"

type AtomicInt32 struct {
	value int32
}

func (self *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&self.value)
}

func (self *AtomicInt32) IncrementAndGet() int32 {
	return atomic.AddInt32(&self.value, 1)
}

func (self *AtomicInt32) GetAndIncrement() (int32) {
	return self.GetAndAdd(1)
}

func (self *AtomicInt32) DecrementAndGet() (int32) {
	return self.AddAndGet(-1)
}

func (self *AtomicInt32) GetAndDecrement() (int32) {
	return self.GetAndAdd(-1)
}

func (self *AtomicInt32) AddAndGet(i int32) (int32) {
	return atomic.AddInt32(&self.value, i)
}

func (self *AtomicInt32) GetAndAdd(i int32) (int32) {
	var ret int32
	for ; ; {
		ret = atomic.LoadInt32(&self.value)
		if atomic.CompareAndSwapInt32(&self.value, ret, ret + i) {
			break
		}
	}
	return ret
}

func (self *AtomicInt32) CompareAndSet(expect int32, update  int32) (bool) {
	return atomic.CompareAndSwapInt32(&self.value, expect, update)
}