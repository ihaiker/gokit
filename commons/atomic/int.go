package atomic

import "sync/atomic"

type AtomicInt struct {
    value int32
}

func (self *AtomicInt) Get() int {
    return int(atomic.LoadInt32(&self.value))
}

func (self *AtomicInt) IncrementAndGet() int {
    return int(atomic.AddInt32(&self.value, 1))
}

func (self *AtomicInt) GetAndIncrement() (int) {
    return self.GetAndAdd(1)
}

func (self *AtomicInt) DecrementAndGet() (int) {
    return self.AddAndGet(-1)
}

func (self *AtomicInt) GetAndDecrement() (int) {
    return self.GetAndAdd(-1)
}

func (self *AtomicInt) AddAndGet(i int) (int) {
    return int(atomic.AddInt32(&self.value, int32(i)))
}

func (self *AtomicInt) GetAndAdd(i int) (int) {
    var ret int32
    for ; ; {
        ret = atomic.LoadInt32(&self.value)
        if atomic.CompareAndSwapInt32(&self.value, ret, ret+int32(i)) {
            break
        }
    }
    return int(ret)
}

func (self *AtomicInt) Set(i int) {
    atomic.StoreInt32(&self.value, int32(i))
}

func (self *AtomicInt) CompareAndSet(expect int, update int) (bool) {
    return atomic.CompareAndSwapInt32(&self.value, int32(expect), int32(update))
}
