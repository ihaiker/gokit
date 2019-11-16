
package atomic

import "sync/atomic"

type AtomicUint64 struct {
	value uint64
}

func NewAtomicUint64(initValue uint64) *AtomicUint64 {
	return &AtomicUint64{value: initValue}
}

func (self *AtomicUint64) Get() (uint64) {
	return atomic.LoadUint64(&self.value)
}

func (self *AtomicUint64) IncrementAndGet(i uint) (uint64) {
	return atomic.AddUint64(&self.value, uint64(i))
}

func (self *AtomicUint64) GetAndIncrement(i uint) (uint64) {
	var ret uint64
	for {
		ret = atomic.LoadUint64(&self.value)
		newValue := ret + uint64(i)
		if atomic.CompareAndSwapUint64(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicUint64) DecrementAndGet(i uint) (uint64) {
	var ret uint64
	for {
		ret = atomic.LoadUint64(&self.value)
		newValue := ret - uint64(i)
		if atomic.CompareAndSwapUint64(&self.value, ret, newValue) {
			return newValue
		}
	}
}

func (self *AtomicUint64) GetAndDecrement(i uint) (uint64) {
	var ret uint64
	for ; ; {
		ret = atomic.LoadUint64(&self.value)
		newValue := ret - uint64(i)
		if atomic.CompareAndSwapUint64(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *AtomicUint64) Set(i uint64) {
	atomic.StoreUint64(&self.value, i)
}

func (self *AtomicUint64) CompareAndSet(expect uint64, update uint64) (bool) {
	return atomic.CompareAndSwapUint64(&self.value, expect, update)
}


