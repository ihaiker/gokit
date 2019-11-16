package lldb

import (
	"github.com/ihaiker/gokit/commons"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

const DEFAULT_CONFIG = `
data:
    path: "./data"
    clean: false

locks:
    keys: 4
    hash: 2
    queue: 2
    set: 2
    sortedset: 2

options:
    compression: none
`

type Config struct {
	Data    struct {
		        //数据存储的位置
		        Path  string
		        //启动时是否清除数据
		        Clean bool
	        }
	//各个数据类型锁的个数（指数），例如定于为1，锁的个数就为2,定义为3锁的个数为8
	Locks   struct {
		        Keys      uint8
		        Hash      uint8
		        Queue     uint8
		        Set       uint8
		        SortedSet uint8
	        }
	Options struct {
		        // default/none/snappy
		        Compression string
	        }
}

func (self *Config) GetDataPath() string {
	if self.Data.Path != "" {
		return self.Data.Path
	}
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA") + string(filepath.Separator) + "lldb" + string(filepath.Separator) + strconv.Itoa(os.Getpid())
	default:
		return os.Getenv("HOME") + string(filepath.Separator) + ".lldb" + string(filepath.Separator) + strconv.Itoa(os.Getpid())
	}
}

func (self *Config) GetOptions() *opt.Options {
	var compression opt.Compression = opt.DefaultCompression
	switch self.Options.Compression {
	case "default":
		compression = opt.DefaultCompression
	case "none":
		compression = opt.NoCompression
	case "snappy":
		compression = opt.SnappyCompression
	}
	return &opt.Options{
		Compression: compression,
	}
}

func (self *Config) Clone() (*Config,error) {
	newConfig := &Config{}
	if err := commons.Copy(newConfig,self); err != nil {
		return nil,err
	}
	return newConfig,nil;
}
