/* 	the local storage collection tools backed by leveldb. */
package lldb

import (
	"os"
	"runtime"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/ihaiker/gokit/files"
	"errors"
	"fmt"
	"log"
	"github.com/syndtr/goleveldb/leveldb/util"
	"path/filepath"
	"io"
)

var LSCL_DEFAULT_PATH string

type LLDBEngine struct {
	data         *leveldb.DB
	meta         *leveldb.DB
	writeOptions *opt.WriteOptions
	readOptions  *opt.ReadOptions
}

func (self *LLDBEngine) Close() error {
	if self.data != nil {
		if err := self.data.Close(); err != nil {
			return err
		}
	}
	if self.meta != nil {
		if err := self.meta.Close(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	switch runtime.GOOS {
	case "windows":
		LSCL_DEFAULT_PATH = os.Getenv("APPDATA") + string(filepath.Separator) + "lldb"
	default:
		LSCL_DEFAULT_PATH = os.Getenv("HOME") + string(filepath.Separator) + ".lldb"
	}
}

//Use the default location initialization `leveldb` library
func New() (*LLDBEngine, error) {
	return NewWithPath(LSCL_DEFAULT_PATH)
}

//Using the location specified initialization `leveldb` library
func NewWithPath(dbPath string) (*LLDBEngine, error) {
	return NewWith(dbPath, nil)
}

func NewWith(path string, opt *opt.Options) (*LLDBEngine, error) {
	log.Printf("init lldb datapath: %s", path)

	if !fileKit.Exist(path) {
		if err := os.MkdirAll(path, os.ModeDir); err != nil {
			return nil, err
		}
	} else if !fileKit.IsDir(path) {
		return nil, errors.New(fmt.Sprintf("the path %s not a folder", path))
	}

	data, err := leveldb.OpenFile(path + "/data", opt)
	if err != nil {
		return nil, err
	}
	meta, err := leveldb.OpenFile(path + "/meta", opt)
	if err != nil {
		return nil, err
	}
	return &LLDBEngine{meta:meta, data:data}, nil
}

func (self *LLDBEngine) Data() *leveldb.DB {
	return self.data
}

func (self *LLDBEngine) Meta() *leveldb.DB {
	return self.meta
}

func (self *LLDBEngine) FlushDB() {
	it := self.data.NewIterator(&util.Range{}, self.readOptions)
	for ; it.Next(); {
		self.data.Delete(it.Key(), self.writeOptions)
	}
}

func (self *LLDBEngine) toTest(out io.Writer) {
	it := self.data.NewIterator(&util.Range{}, self.readOptions)
	for ; it.Next(); {
		fmt.Fprint(out, string(it.Key())," = ", string(it.Value()),"\n")
	}
}