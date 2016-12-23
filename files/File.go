package fileKit

import (
	"os"
	"errors"
	"path/filepath"
	"io/ioutil"
	"github.com/ihaiker/gokit/commons"
)

type File struct {
	path string
}

func New(path string) *File {
	absPath, _ := filepath.Abs(path)
	return &File{path:absPath}
}

func (self *File) IsFile() bool {
	return !IsDir(self.path)
}
func (self *File) IsDir() bool {
	return IsDir(self.path)
}
func (self *File) Exist() bool {
	return Exist(self.path)
}
func (self *File) Name() string {
	return filepath.Base(self.path)
}

func (self *File) Mkdir() error {
	if self.Exist() {
		return errors.New("the fodler or file is exits")
	}
	return os.Mkdir(self.path, os.ModePerm)
}

func (self *File) Parent() *File {
	if filepath.Dir(self.path) == self.path {
		return nil
	}
	return New(filepath.Dir(self.path))
}

func (self *File) Equal(file *File) bool {
	return self.path == file.path
}

func (self *File) List() ([]*File, error) {
	dir, err := ioutil.ReadDir(self.path)
	if err != nil {
		return nil, err
	}
	files := make([]*File, len(dir))
	for idx, finfo := range dir {
		files[idx] = New(self.path + "/" + finfo.Name())
	}
	return files, nil
}
func (self *File) ToString() string {
	return self.path
}

func (self *File) ToBytes() ([]byte, error) {
	return ioutil.ReadFile(self.path)
}

func (self *File) GetWriter(append bool) (*os.File, error) {
	flag := os.O_RDWR | os.O_CREATE
	if append && self.Exist() {
		flag = flag | os.O_APPEND
	}
	return os.OpenFile(self.path, flag, 0666)
}

func (self *File) LineIterator() (commonKit.CloseIterator, error) {
	if self.Exist() && self.IsFile() {
		f, err := os.OpenFile(self.path, (os.O_RDWR | os.O_APPEND), 0666)
		if err != nil {
			return nil, err
		} else {
			return newIterator(f), nil
		}
	}
	return nil, errors.New("not found or is not file")
}
