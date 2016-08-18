package json

import (
	"errors"
	"io/ioutil"
	"strings"
	"encoding/json"
	"reflect"
)

//注册值的替换方式
type JsonValueReplace interface {
	Replace(kind reflect.Kind, path string, src interface{}, dst interface{}) (interface{}, error)
}

type JsonMerger struct {
	src     []byte                            //需要被扩展的json
	dst     []byte                            //扩展内容
	out     []byte                            //最终结果
	methods map[reflect.Kind]JsonValueReplace //值的替换方式
}

//设置需要扩展的json内容字节
func (self *JsonMerger) SetSrcByte(src []byte) error {
	if len(src) == 0 {
		return errors.New("src is null")
	}
	self.src = src
	return nil
}

//设置需要扩展的json内容
func (self *JsonMerger) SetSrcString(src string) error {
	if src == "" || strings.Compare(src, "") == 0 {
		return errors.New("src is null")
	}
	self.SetSrcByte(([]byte)(src))
	return nil
}

//设置需要扩展内容所在文件
func (self *JsonMerger) SetSrcFile(srcFile string) error {
	body, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	self.SetSrcByte(body)
	return nil
}

func (self *JsonMerger) SetDstByte(dst []byte) {
	self.dst = dst
}

func (self *JsonMerger) SetDstString(dst string) error {
	if dst == "" {
		return errors.New("src is null")
	}
	self.SetDstByte(([]byte)(dst))
	return nil
}

func (self *JsonMerger) SetDstFile(srcFile string) error {
	body, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	self.SetDstByte(body)
	return nil
}

func (self *JsonMerger) ToString() string {
	if self.out == nil {
		self.Merge()
	}
	return string(self.out)
}

// 返回最终结果
func (self *JsonMerger) GetOut() ([]byte, error) {
	if self.out == nil {
		if err := self.Merge(); err != nil {
			return nil, err
		}
	}
	return self.out, nil
}

//根据out产生一个新的src并且情况dst,out
func (self *JsonMerger) NewFromOut() error {
	if self.out == nil {
		if err := self.Merge(); err != nil {
			return err
		}
	}
	self.dst = ([]byte)(nil)
	self.SetSrcByte(self.out)
	self.out = ([]byte)(nil)
	return nil
}

func merge(srcJson, dstJson map[string]interface{}, path string) (map[string]interface{}, error) {
	outJson := srcJson
	for key, value := range dstJson {
		current_path := path + "." + key
		srcValue := srcJson[key]
		typeof := reflect.TypeOf(value)

		switch typeof.Kind() {
		case reflect.Struct       : fallthrough
		case reflect.UnsafePointer: fallthrough
		case reflect.Ptr          : fallthrough
		case reflect.Chan         : fallthrough
		case reflect.Func         : fallthrough
		case reflect.Interface    : fallthrough
		case reflect.Complex64    : fallthrough
		case reflect.Complex128   : fallthrough
		case reflect.Uintptr      : fallthrough
		case reflect.Invalid      : return nil, errors.New("unsupport " + typeof.String() + " in (" + current_path + ")")
		case reflect.Bool         : fallthrough
		case reflect.Int          : fallthrough
		case reflect.Int8         : fallthrough
		case reflect.Int16        : fallthrough
		case reflect.Int32        : fallthrough
		case reflect.Int64        : fallthrough
		case reflect.Uint         : fallthrough
		case reflect.Uint8        : fallthrough
		case reflect.Uint16       : fallthrough
		case reflect.Uint32       : fallthrough
		case reflect.Uint64       : fallthrough
		case reflect.Float32      : fallthrough
		case reflect.Float64      : fallthrough
		case reflect.String       : outJson[key] = value
		case reflect.Array        : fallthrough
		case reflect.Slice        : outJson[key] = value
		case reflect.Map          :
			if srcValue == nil {
				srcValue = make(map[string]interface{})
			}
			outValue, err := merge(srcValue.(map[string]interface{}), value.(map[string]interface{}), current_path)
			if err != nil {
				return nil, err
			}
			outJson[key] = outValue
		}
	}
	return outJson, nil
}

func (self *JsonMerger) Merge() error {
	if self.out != nil {
		return nil
	}

	if self.src == nil || len(self.src) == 0 {
		return errors.New("src is null")
	}
	if self.dst == nil || len(self.dst) == 0 {
		return errors.New("dst is null")
	}
	var srcJson, dstJson map[string]interface{}

	if err := json.Unmarshal(self.src, &srcJson); err != nil {
		return err
	}
	if err := json.Unmarshal(self.dst, &dstJson); err != nil {
		return err
	}

	outJson, err := merge(srcJson, dstJson, "")
	if err != nil {
		return err
	}

	out, err := json.Marshal(outJson)
	if err != nil {
		return err
	}
	self.out = out
	return nil
}