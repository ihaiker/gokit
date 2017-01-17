package main

import (
	"fmt"
	"github.com/ihaiker/gokit/protocol/redis"
	"github.com/ihaiker/gokit/lldb/server/lldbs"
	"github.com/ihaiker/gokit/lldb"
)

func ifPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf("Panic: %v\n", msg)
		}
	}()
	//load config
	conf, err := lldb.SetConfig(""); ifPanic(err)
	redisHandler := lldbs.NewRedisHandler(conf); ifPanic(err)
	err = redisHandler.Select(0); ifPanic(err)
	srv, err := redis.NewServer(redis.DefaultConfig().Handler(redisHandler)); ifPanic(err)
	err = srv.ListenAndServe(); ifPanic(err)
}
