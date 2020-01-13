package atomic

import "sync/atomic"

type AtomicInt struct {
	value int32
}

func NewAtomicInt(initValue int) *AtomicInt {
	return &AtomicInt{value: int32(initValue)}
}

func (self *AtomicInt) Get() (int) {
	return int(atomic.LoadInt32(&self.value))
}

func (self *AtomicInt) IncrementAndGet(i uint) (int) {
	return int(atomic.AddInt32(&self.value, int32(i)))
}

func (self *AtomicInt) GetAndIncrement(i uint) (int) {
	var ret int32
	for {
		ret = atomic.LoadInt32(&self.value)
		newValue := ret + int32(i)
		if atomic.CompareAndSwapInt32(&self.value, ret, newValue) {
			return int(ret)
		}
	}
}

func (self *AtomicInt) DecrementAndGet(i uint) (int) {
	var ret int32
	for {
		ret = atomic.LoadInt32(&self.value)
		newValue := ret - int32(i)
		if atomic.CompareAndSwapInt32(&self.value, ret, newValue) {
			return int(newValue)
		}
	}
}

func (self *AtomicInt) GetAndDecrement(i uint) (int) {
	var ret int32
	for ; ; {
		ret = atomic.LoadInt32(&self.value)
		newValue := ret - int32(i)
		if atomic.CompareAndSwapInt32(&self.value, ret, newValue) {
			return int(ret)
		}
	}
}

func (self *AtomicInt) Set(i int) {
	atomic.StoreInt32(&self.value, int32(i))
}

func (self *AtomicInt) CompareAndSet(expect int, update int) (bool) {
	return atomic.CompareAndSwapInt32(&self.value, int32(expect), int32(update))
}
