
package atomic

import "sync/atomic"

type AtomicInt64 struct {
	value int64
}

func NewAtomicInt64(initValue int64) *AtomicInt64 {
	return &AtomicInt64{value: initValue}
}

func (self *AtomicInt64) Get() (int64) {
	return atomic.LoadInt64(&self.value)
}

func (self *AtomicInt64) IncrementAndGet(i uint) (int64) {
	return atomic.AddInt64(&self.value, int64(i))
}

func (self *AtomicInt64) GetAndIncrement(i uint) (int64) {
	var ret int64
	for {
		ret = atomic.LoadInt64(&self.value)
		newValue := ret + int64(i)
		if atomic.CompareAndSwapInt64(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicInt64) DecrementAndGet(i uint) (int64) {
	var ret int64
	for {
		ret = atomic.LoadInt64(&self.value)
		newValue := ret - int64(i)
		if atomic.CompareAndSwapInt64(&self.value, ret, newValue) {
			return newValue
		}
	}
}

func (self *AtomicInt64) GetAndDecrement(i uint) (int64) {
	var ret int64
	for ; ; {
		ret = atomic.LoadInt64(&self.value)
		newValue := ret - int64(i)
		if atomic.CompareAndSwapInt64(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicInt64) Set(i int64) {
	atomic.StoreInt64(&self.value, i)
}

func (self *AtomicInt64) CompareAndSet(expect int64, update int64) (bool) {
	return atomic.CompareAndSwapInt64(&self.value, expect, update)
}


