package lldb

import (
	"sync"
	"github.com/bluele/gcache"
	"time"
)

var cache gcache.Cache

func init() {
	cache = gcache.New(10).LRU().
		Expiration(time.Hour).LoaderFunc(func(key interface{}) (interface{}, error) {
		return &sync.RWMutex{}, nil
	}).Build()
}

func NewRowLock(key string) (*sync.RWMutex, error) {
	lock, err := cache.Get(key)
	if err != nil {
		return nil, err
	}
	return lock.(*sync.RWMutex), err
}
