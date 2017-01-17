/*
    本实现借鉴了 github.com/dotcloud/go-redis-server 其实就是完全的抄袭
*/
package main

import (
	"fmt"
	"github.com/ihaiker/gokit/protocol/redis"
	"github.com/ihaiker/gokit/protocol/redis/example/core"
)

func main() {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf("Panic: %v\n", msg)
		}
	}()
	srv := redis.NewServer()
	srv.RegisterHandler(&core.DefaultHandler{})
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
