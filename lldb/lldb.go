/* 	the local storage collection tools backed by leveldb. */
package lldb

import (
    "os"
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
    "github.com/ihaiker/gokit/files"
    "errors"
    "fmt"
    "log"
    "github.com/syndtr/goleveldb/leveldb/util"
    "bytes"
)

type LLDBEngine struct {
    data         *leveldb.DB
    writeOptions *opt.WriteOptions
    readOptions  *opt.ReadOptions
    config       *Config
}

//close the leveldb connect.
func (self *LLDBEngine) Close() error {
    if self.data != nil {
        if err := self.data.Close(); err != nil {
            return err
        }
    }
    return nil
}

func (self *LLDBEngine) FlushDB() error {
    sn, err := self.data.GetSnapshot()
    if err != nil {
        return err
    }
    it := sn.NewIterator(&util.Range{}, self.readOptions)
    for ; it.Next(); {
        if err := self.data.Delete(it.Key(), self.writeOptions); err != nil {
            return err
        }
    }
    return nil
}

func (self *LLDBEngine) toTest() string {
    it := self.data.NewIterator(&util.Range{}, self.readOptions)
    w := bytes.NewBufferString("")
    for ; it.Next(); {
        w.WriteString(string(it.Key()))
        w.WriteString(" = ")
        w.WriteString(string(it.Value()))
        w.WriteString("\n")
    }
    return w.String()
}

//Use the default location initialization `leveldb` library
func Default() (*LLDBEngine, error) {
    cfg,err := SetConfig("")
    if err != nil {
        return nil,err
    }
    return New(cfg)
}

//Using the location specified initialization `leveldb` library
func NewWith(cfgPath string) (*LLDBEngine, error) {
    if !fileKit.IsExistFile(cfgPath) {
        return nil,errors.New("the config file not found !")
    }
    cfg,err := SetConfig(cfgPath)
    if err != nil {
        return nil,err
    }
    return New(cfg)
}

func New(cfg *Config) (*LLDBEngine, error) {
    path := cfg.GetDataPath()
    log.Printf("init lldb datapath: %s", path)
    
    if !fileKit.Exist(path) {
        if err := os.MkdirAll(path, os.ModeDir); err != nil {
            return nil, err
        }
    } else if !fileKit.IsDir(path) {
        return nil, errors.New(fmt.Sprintf("the path %s not a folder", path))
    }
    data, err := leveldb.OpenFile(path, cfg.GetOptions())
    if err != nil {
        return nil, err
    }
    
    return &LLDBEngine{data:data,config:cfg}, nil
}

