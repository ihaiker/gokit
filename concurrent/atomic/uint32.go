
package atomic

import "sync/atomic"

type AtomicUint32 struct {
	value uint32
}

func NewAtomicUint32(initValue uint32) *AtomicUint32 {
	return &AtomicUint32{value: initValue}
}

func (self *AtomicUint32) Get() (uint32) {
	return atomic.LoadUint32(&self.value)
}

func (self *AtomicUint32) IncrementAndGet(i uint) (uint32) {
	return atomic.AddUint32(&self.value, uint32(i))
}

func (self *AtomicUint32) GetAndIncrement(i uint) (uint32) {
	var ret uint32
	for {
		ret = atomic.LoadUint32(&self.value)
		newValue := ret + uint32(i)
		if atomic.CompareAndSwapUint32(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicUint32) DecrementAndGet(i uint) (uint32) {
	var ret uint32
	for {
		ret = atomic.LoadUint32(&self.value)
		newValue := ret - uint32(i)
		if atomic.CompareAndSwapUint32(&self.value, ret, newValue) {
			return newValue
		}
	}
}

func (self *AtomicUint32) GetAndDecrement(i uint) (uint32) {
	var ret uint32
	for ; ; {
		ret = atomic.LoadUint32(&self.value)
		newValue := ret - uint32(i)
		if atomic.CompareAndSwapUint32(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicUint32) Set(i uint32) {
	atomic.StoreUint32(&self.value, i)
}

func (self *AtomicUint32) CompareAndSet(expect uint32, update uint32) (bool) {
	return atomic.CompareAndSwapUint32(&self.value, expect, update)
}


