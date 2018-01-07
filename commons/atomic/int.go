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

func (self *AtomicInt32) Set(i int32) {
    atomic.StoreInt32(&self.value,i)
}

func (self *AtomicInt32) CompareAndSet(expect int32, update  int32) (bool) {
	return atomic.CompareAndSwapInt32(&self.value, expect, update)
}




type AtomicUInt32 struct {
	value uint32
}

func (self *AtomicUInt32) Get() uint32 {
	return atomic.LoadUint32(&self.value)
}

func (self *AtomicUInt32) IncrementAndGet() uint32 {
	return atomic.AddUint32(&self.value, 1)
}

func (self *AtomicUInt32) GetAndIncrement() (uint32) {
	return self.GetAndAdd(1)
}

func (self *AtomicUInt32) DecrementAndGet() (uint32) {
	var old uint32
	for ; ; {
		old = atomic.LoadUint32(&self.value)
		if atomic.CompareAndSwapUint32(&self.value, old, old - 1 ) {
			break
		}
	}
	return old - 1
}

func (self *AtomicUInt32) GetAndDecrement() (uint32) {
	var old uint32
	for ; ; {
		old = atomic.LoadUint32(&self.value)
		if atomic.CompareAndSwapUint32(&self.value, old, old - 1 ) {
			break
		}
	}
	return old
}

func (self *AtomicUInt32) AddAndGet(i uint32) (uint32) {
	return atomic.AddUint32(&self.value, i)
}

func (self *AtomicUInt32) GetAndAdd(i uint32) (uint32) {
	var ret uint32
	for ; ; {
		ret = atomic.LoadUint32(&self.value)
		if atomic.CompareAndSwapUint32(&self.value, ret, ret + i) {
			break
		}
	}
	return ret
}

func (self *AtomicUInt32) Set(i uint32) {
    atomic.StoreUint32(&self.value,i)
}

func (self *AtomicUInt32) CompareAndSet(expect uint32, update  uint32) (bool) {
	return atomic.CompareAndSwapUint32(&self.value, expect, update)
}