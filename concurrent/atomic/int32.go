
package atomic

import "sync/atomic"

type AtomicInt32 struct {
	value int32
}

func NewAtomicInt32(initValue int32) *AtomicInt32 {
	return &AtomicInt32{value: initValue}
}

func (self *AtomicInt32) Get() (int32) {
	return atomic.LoadInt32(&self.value)
}

func (self *AtomicInt32) IncrementAndGet(i uint) (int32) {
	return atomic.AddInt32(&self.value, int32(i))
}

func (self *AtomicInt32) GetAndIncrement(i uint) (int32) {
	var ret int32
	for {
		ret = atomic.LoadInt32(&self.value)
		newValue := ret + int32(i)
		if atomic.CompareAndSwapInt32(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicInt32) DecrementAndGet(i uint) (int32) {
	var ret int32
	for {
		ret = atomic.LoadInt32(&self.value)
		newValue := ret - int32(i)
		if atomic.CompareAndSwapInt32(&self.value, ret, newValue) {
			return newValue
		}
	}
}

func (self *AtomicInt32) GetAndDecrement(i uint) (int32) {
	var ret int32
	for ; ; {
		ret = atomic.LoadInt32(&self.value)
		newValue := ret - int32(i)
		if atomic.CompareAndSwapInt32(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicInt32) Set(i int32) {
	atomic.StoreInt32(&self.value, i)
}

func (self *AtomicInt32) CompareAndSet(expect int32, update int32) (bool) {
	return atomic.CompareAndSwapInt32(&self.value, expect, update)
}


