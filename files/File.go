package fileKit

import (
    "os"
    "errors"
    "path/filepath"
    "io/ioutil"
    "github.com/ihaiker/gokit/commons"
    "github.com/hpcloud/tail"
)

type File struct {
    path string
}

func New(path string) *File {
    absPath, _ := filepath.Abs(path)
    return &File{path: absPath}
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

func (self *File) Rename(newName string) error {
    dir, _ := filepath.Split(self.path)
    newPath := dir + "/" + newName
    return os.Rename(self.path, newPath)
}

func (self *File) GetPath() string {
    return self.path
}

func (self *File) ToString() (string, error) {
    bs, err := self.ToBytes()
    return string(bs), err
}

//delete file or folder
func (self *File) Remove() error {
    return os.Remove(self.path)
}

//delete file or folder and subfolder
func (self *File) RemoveAll() error {
    return os.RemoveAll(self.path)
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

func (self *File) Size() int64 {
    if self.IsFile() {
        f, _ := os.Stat(self.path);
        return f.Size()
    } else {
        return -1
    }
}

func (self *File) LineIterator() (commonKit.CloseIterator, error) {
    if self.Exist() && self.IsFile() {
        f, err := self.GetReader()
        if err != nil {
            return nil, err
        } else {
            return newIterator(f), nil
        }
    }
    return nil, errors.New("not found or is not file")
}

func (self *File) Tail() (*tail.Tail, error) {
    if self.Exist() && self.IsFile() {
        return tail.TailFile(self.GetPath(), tail.Config{Follow: true, MustExist: true})
    }
    return nil, errors.New("not found or is not file")
}

func (self *File) GetReader() (*os.File, error) {
    if self.Exist() && self.IsFile() {
        return os.OpenFile(self.path, (os.O_RDWR | os.O_APPEND), 0666)
    }
    return nil, errors.New("not found or is not file")
}
