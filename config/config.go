package config

import (
	"github.com/ihaiker/gokit/files"
	"reflect"
	"errors"
	"github.com/ihaiker/gokit/convert"
	"strings"
)

var NOT_FOUND = errors.New("not found")

type Unmarshal func([]byte, interface{}) error

type Config struct {
	data      []byte
	unmarshal Unmarshal
	merge     *Merger
}

func (self *Config) Unmarshal(v interface{}) error {
	return self.unmarshal(self.data, v)
}

func (self *Config) UnmarshalP(v interface{}) {
	if err := self.Unmarshal(v); err != nil {
		panic(err)
	}
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
		} else {
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

func (self *Config) LoadP(items... interface{}) {
	if err := self.Load(items...); err != nil {
		panic(err)
	}
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

func (self *Config) convertMap(v interface{}) (map[string]interface{}, error) {
	//兼容yaml、json两种格式问题
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{}), nil
	case map[interface{}]interface{}:
		newObj := map[string]interface{}{}
		for k, v := range v.(map[interface{}]interface{}) {
			newObj[convertKit.SafeString(k)] = v
		}
		return newObj, nil
	default:
		return nil, errors.New("Failed to find child nodes")
	}
}

func (self *Config) get(path string) (interface{}, error) {
	var obj map[string]interface{}
	if err := self.Unmarshal(&obj); err != nil {
		return nil, err
	}
	paths := strings.Split(path, ".")
	for index, p := range paths {
		if data, has := obj[p]; has {
			if len(paths) - 1 == index {
				return data, nil
			} else {
				if m, err := self.convertMap(data); err != nil {
					return nil, err
				} else {
					obj = m
				}
			}
		} else {
			return nil, NOT_FOUND
		}
	}
	return nil, errors.New("Unkown error!")
}

func (self *Config) GetString(path string) (string, error) {
	if i, err := self.get(path); err != nil {
		return "", err
	} else {
		return convertKit.String(i)
	}
}

func (self *Config) GetStringP(path string) string {
	if str, err := self.GetString(path); err != nil {
		panic(err)
	} else {
		return str
	}
}

func (self *Config) Get(path string) (map[string]interface{}, error) {
	if v, err := self.get(path); err != nil {
		return nil, err
	} else {
		return self.convertMap(v)
	}
}

func (self *Config) GetP(path string) map[string]interface{} {
	if m, err := self.Get(path); err != nil {
		panic(err)
	} else {
		return m
	}
}

func (self *Config) GetInt(path string) (int, error) {
	if v, err := self.get(path); err != nil {
		return 0, err
	} else {
		return convertKit.Int(v)
	}
}

func (self *Config) GetIntP(path string) int {
	if i, err := self.GetInt(path); err != nil {
		panic(err)
	} else {
		return i
	}
}

func (self *Config) GetSlice(path string) ([]interface{}, error) {
	if i, err := self.get(path); err != nil {
		return nil, err
	} else {
		switch i.(type) {
		case []interface{}:
			return i.([]interface{}), nil
		default:
			return nil, errors.New("not slice in path: " + path)
		}
	}
}

func (self *Config) GetSliceP(path string) []interface{} {
	if s, err := self.GetSlice(path); err != nil {
		panic(err)
	} else {
		return s
	}
}

func NewConfig(def []byte, unmarshal Unmarshal, merge *Merger) *Config {
	return &Config{
		data:def,
		unmarshal:unmarshal,
		merge:merge,
	}
}