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





type AtomicUInt64 struct {
	value uint64
}

func (self *AtomicUInt64) Get() uint64 {
	return atomic.LoadUint64(&self.value)
}

func (self *AtomicUInt64) IncrementAndGet() uint64 {
	return atomic.AddUint64(&self.value, 1)
}

func (self *AtomicUInt64) GetAndIncrement() (uint64) {
	return self.GetAndAdd(1)
}


func (self *AtomicUInt64) DecrementAndGet() (uint64) {
	var old uint64
	for ; ; {
		old = atomic.LoadUint64(&self.value)
		if atomic.CompareAndSwapUint64(&self.value, old, old - 1 ) {
			break
		}
	}
	return old - 1
}

func (self *AtomicUInt64) GetAndDecrement() (uint64) {
	var old uint64
	for ; ; {
		old = atomic.LoadUint64(&self.value)
		if atomic.CompareAndSwapUint64(&self.value, old, old - 1 ) {
			break
		}
	}
	return old
}

func (self *AtomicUInt64) AddAndGet(i uint64) (uint64) {
	return atomic.AddUint64(&self.value, i)
}

func (self *AtomicUInt64) GetAndAdd(i uint64) (uint64) {
	var ret uint64
	for ; ; {
		ret = atomic.LoadUint64(&self.value)
		if atomic.CompareAndSwapUint64(&self.value, ret, ret + i) {
			break
		}
	}
	return ret
}

func (self *AtomicUInt64) CompareAndSet(expect uint64, update  uint64) (bool) {
	return atomic.CompareAndSwapUint64(&self.value, expect, update)
}