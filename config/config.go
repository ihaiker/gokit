package config

import (
	"github.com/ihaiker/gokit/files"
	"reflect"
	"errors"
)

type Unmarshal func ([]byte,interface{}) error

type Config struct {
	data []byte
	unmarshal Unmarshal
	merge *Merger
}

func (self *Config) Unmarshal(v interface{}) error {
	return self.unmarshal(self.data, v)
}

func (self *Config) ToString() string {
	return string(self.data)
}

func (self *Config) Load(items ...interface{}) error {
	if len(items) == 0 {
		return errors.New("the length is 0")
	}

	if err := self.merge.SetSrcByte(self.data); err != nil {
		return err
	}
	for idx, item := range items {
		if idx != 0 {
			if err := self.merge.NewFromOut(); err != nil {
				return err
			}
		}
		if out, err := self.getJsonContentByte(item); err != nil {
			return err
		}else{
			self.merge.SetDstByte(out)
		}
	}

	if err := self.merge.Merge(); err != nil {
		return err
	}
	out, err := self.merge.GetOut()
	if err != nil {
		return err
	}
	self.data = out
	return nil
}

//获取内容体字节
func (self *Config) getJsonContentByte(item interface{}) ([]byte, error) {
	switch item.(type) {
	case string:
		return []byte(item.(string)), nil
	case *fileKit.File:
		return item.(*fileKit.File).ToBytes()
	default:
		return nil, errors.New("unsupport type" + reflect.TypeOf(item).Name())
	}
}

func NewConfig(def []byte,unmarshal Unmarshal,merge *Merger) *Config{
	return &Config{
		data:def,
		unmarshal:unmarshal,
		merge:merge,
	}
}