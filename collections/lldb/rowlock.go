package lldb

import (
	"sync"
	"math"
	"errors"
)

type Locks struct {
	d map[uint8]*sync.RWMutex
	y uint8
}

func (self *Locks) Get(key string) *sync.RWMutex {
	n := uint8(key[0])
	return self.d[n%self.y]
}

func NewLocks(num uint8) (*Locks, error) {
	if num > 8 {
		return nil, errors.New("the numbers must be smaller than 8")
	}
	n := uint8(math.Pow(2, float64(num)))
	locks := &Locks{y:n}
	locks.d = make(map[uint8]*sync.RWMutex, n)
	for i := uint8(0); i < n; i++ {
		locks.d[i] = &sync.RWMutex{}
	}
	return locks, nil
}