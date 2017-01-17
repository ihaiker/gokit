package lldb

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/config/yaml"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/ihaiker/gokit/commons/logs"
	"github.com/ihaiker/gokit/commons"
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

logs:
    level: debug
    path: "./logs/server.log"
    
`

type Config struct {
	Logs    struct {
		        Level string
		        Path  string
	        }
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
	if err := commonKit.Copy(newConfig,self); err != nil {
		return nil,err
	}
	return newConfig,nil;
}

//read the config file
//参数 file 可以为空，也可以是一个未找到的文件
func SetConfig(file string) (*Config, error) {
	configTools, err := yaml.Config(DEFAULT_CONFIG)
	if err != nil {
		return nil, err
	}
	
	if fileKit.IsExistFile(file) {
		logs.Info("load config file ", file)
		err = configTools.Load(fileKit.New(file))
		if err != nil {
			return nil, err
		}
	}
	cfg := &Config{}
	if err = configTools.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}