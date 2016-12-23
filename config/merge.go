package config

import (
	"strings"
	"io/ioutil"
	"errors"
)

type merger  func (src []byte,dst []byte,path string) (out []byte,err error)

type Merger struct {
	src []byte //需要被扩展的json
	dst []byte //扩展内容
	out []byte //最终结果
	merger
}

//设置需要扩展的json内容字节
func (self *Merger) SetSrcByte(src []byte) error {
	if len(src) == 0 {
		return errors.New("src is null")
	}
	self.src = src
	self.out = nil
	return nil
}

//设置需要扩展的json内容
func (self *Merger) SetSrcString(src string) error {
	if src == "" || strings.Compare(src, "") == 0 {
		return errors.New("src is null")
	}
	self.SetSrcByte(([]byte)(src))
	return nil
}

//设置需要扩展内容所在文件
func (self *Merger) SetSrcFile(srcFile string) error {
	body, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	self.SetSrcByte(body)
	return nil
}

func (self *Merger) SetDstByte(dst []byte) {
	self.dst = dst
	self.out = nil
}

func (self *Merger) SetDstString(dst string) error {
	if dst == "" {
		return errors.New("src is null")
	}
	self.SetDstByte(([]byte)(dst))
	return nil
}

func (self *Merger) SetDstFile(srcFile string) error {
	body, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	self.SetDstByte(body)
	return nil
}

func (self *Merger) ToString() string {
	if self.out == nil {
		self.Merge()
	}
	return string(self.out)
}
func (self *Merger) GetOut() ([]byte, error) {
	if err := self.Merge(); err != nil {
		return nil, err
	}
	return self.out, nil
}
func (self *Merger) NewFromOut() error {
	if err := self.Merge(); err != nil {
		return err
	}
	self.dst = ([]byte)(nil)
	self.SetSrcByte(self.out)
	self.out = ([]byte)(nil)
	return nil
}


func (self *Merger) Merge() error {
	if self.out != nil {
		return nil
	}

	if self.src == nil || len(self.src) == 0 {
		return errors.New("src is null")
	}
	if self.dst == nil || len(self.dst) == 0 {
		return errors.New("dst is null")
	}
	out,err := self.merger(self.src,self.dst,"")

	if err != nil {
		return err
	}
	self.out = out
	return nil
}

func NewMerger(merger merger) *Merger{
	return &Merger{merger:merger}
}