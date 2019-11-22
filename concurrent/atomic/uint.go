package atomic

import "sync/atomic"

type AtomicUint struct {
	value uint32
}

func NewAtomicUint(initValue int) *AtomicUint {
	return &AtomicUint{value: uint32(initValue)}
}

func (self *AtomicUint) Get() (uint) {
	return uint(atomic.LoadUint32(&self.value))
}

func (self *AtomicUint) IncrementAndGet(i uint) (uint) {
	return uint(atomic.AddUint32(&self.value, uint32(i)))
}

func (self *AtomicUint) GetAndIncrement(i uint) (uint) {
	var ret uint32
	for {
		ret = atomic.LoadUint32(&self.value)
		newValue := ret + uint32(i)
		if atomic.CompareAndSwapUint32(&self.value, ret, newValue) {
			return uint(ret)
		}
	}
}

func (self *AtomicUint) DecrementAndGet(i uint) (uint) {
	var ret uint32
	for {
		ret = atomic.LoadUint32(&self.value)
		newValue := ret - uint32(i)
		if atomic.CompareAndSwapUint32(&self.value, ret, newValue) {
			return uint(newValue)
		}
	}
}

func (self *AtomicUint) GetAndDecrement(i uint) (uint) {
	var ret uint32
	for ; ; {
		ret = atomic.LoadUint32(&self.value)
		newValue := ret - uint32(i)
		if atomic.CompareAndSwapUint32(&self.value, ret, newValue) {
			return uint(ret)
		}
	}
}

func (self *AtomicUint) Set(i int) {
	atomic.StoreUint32(&self.value, uint32(i))
}

func (self *AtomicUint) CompareAndSet(expect int, update int) (bool) {
	return atomic.CompareAndSwapUint32(&self.value, uint32(expect), uint32(expect))
}
