package lldbs

import (
)
import "github.com/ihaiker/gokit/lldb"

//所有redis的协议内容
type redisHandler struct {
	db []*lldb.LLDBEngine
}
//key
func (self *redisHandler) Set(key string, value []byte) error {
	return nil
}
func (self *redisHandler) Del(keys []string) (int, error) {
	return 0, nil
}
func (self *redisHandler) Get(key string) ([]byte, error) {
	return nil,nil
}
func (self *redisHandler) Select(db int) error {
	return nil
}

func NewRedisHandler(conf *lldb.Config) *redisHandler {
	db := make([]*lldb.LLDBEngine,12)
	handler := &redisHandler{db:db}
	return handler
}