/*
    本实现借鉴了 github.com/dotcloud/go-redis-server 其实就是完全的抄袭
*/
package main

import (
	"fmt"
	"github.com/ihaiker/gokit/protocol/redis"
)

type MyHandler struct {
	redis.DefaultHandler
}

//func (h *MyHandler) Set(key string,value []byte) (int,error) {
//    err := h.DefaultHandler.Set(key,value)
//    return 1,err
//}

func (h *MyHandler) Get(key string) ([]byte, error) {
	ret, err := h.DefaultHandler.Get(key)
	if ret == nil {
		return nil, err
	}
	return []byte("BEAM/" + string(ret)), err
}

func main() {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf("Panic: %v\n", msg)
		}
	}()
	srv, err := redis.NewServer(redis.DefaultConfig().Handler(&MyHandler{}))
	if err != nil {
		panic(err)
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
