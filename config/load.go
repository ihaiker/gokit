package config

import (
	"errors"
	"reflect"
	mj "github.com/ihaiker/gokit/json"
	"io/ioutil"
	"encoding/json"
)

type ConfigFile string

type Config struct {
	json_data []byte
}

func (self *Config) Unmarshal(v interface{}) error {
	return json.Unmarshal(self.json_data, v)
}

func getJsonContentByte(item interface{}) ([]byte, error) {
	switch item.(type) {
	case string:
		return []byte(item.(string)), nil
	case ConfigFile:
		return ioutil.ReadFile((string)(item.(ConfigFile)))
	default:
		return nil, errors.New("unsupport type" + reflect.TypeOf(item).Name())
	}
}

//加载所有配置文件
func Load(def interface{}, items... interface{}) (Config, error) {
	cfg := Config{json_data:nil}

	//先把默认的配置加载上来
	out, err := getJsonContentByte(def)
	if err != nil {
		return cfg, err
	}

	if len(items) == 0 {
		//如果不包含数字直接返回
		cfg.json_data = out
		return cfg, err
	}else {
		merge := mj.JsonMerger{}
		if err = merge.SetSrcByte(out); err != nil {
			return cfg, err
		}

		for idx, item := range items {
			if idx != 0 {
				if err = merge.NewFromOut(); err != nil {
					return cfg, err
				}
			}
			if out, err = getJsonContentByte(item); err != nil {
				return cfg, err
			}
			merge.SetDstByte(out)
		}

		if err = merge.Merge(); err != nil {
			return cfg, err
		}
		out, err = merge.GetOut()
		if err != nil {
			return cfg, err
		}
		cfg.json_data = out
		return cfg, err
	}
}