package appenders

import (
	"bufio"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/ihaiker/gokit/files"
	"path/filepath"
	"os"
	"time"
	"errors"
	"io"
	"strings"
	"regexp"
)

//每日回滚日志
type dailyRollingFile struct {
	cache  gcache.Cache
	layout string
}

func (self *dailyRollingFile) Write(p []byte) (n int, err error) {
	key := time.Now().Format(self.layout)
	if fd, err := self.cache.Get(key); err != nil {
		return 0, err
	} else {
		w := fd.(io.Writer)
		return w.Write(p)
	}
}

func (self *dailyRollingFile) Close() error {
	for _, out := range self.cache.GetALL(false) {
		_ = out.(*os.File).Close()
	}
	return nil
}

//创建每日回滚日志的文件夹
func _create_file_dir(logDir string) error {
	if files.NotExist(logDir) {
		if err := os.MkdirAll(logDir, os.ModeDir); err != nil {
			return errors.New("can mkdir " + logDir + " " + err.Error())
		}
	}
	return nil
}

func MatchDailyRollingFile(fileName string) (layout string, match bool) {
	reg, _ := regexp.Compile("\\{[0-9:-]*\\}")
	if !reg.MatchString(fileName) {
		return "", false
	}
	layout = reg.FindString(fileName)
	layout = layout[1 : len(layout)-1]
	return layout, true
}

func NewDailyRollingFileOut(fileName string) (io.Writer, error) {
	logDir := filepath.Dir(fileName)
	if err := _create_file_dir(logDir); err != nil {
		return nil, err
	}
	layout, _ := MatchDailyRollingFile(fileName)

	cb := gcache.New(1).LRU().Expiration(time.Hour)
	cb.LoaderFunc(func(key interface{}) (interface{}, error) {
		fileName := strings.Replace(fileName, "{"+layout+"}", key.(string), 1)
		return newFileout(fileName)
	}).EvictedFunc(func(key, value interface{}) {
		fd := value.(*os.File)
		fmt.Println("close: ", key)
		_ = fd.Close()
	})

	return &dailyRollingFile{
		cache:  cb.Build(),
		layout: layout,
	}, nil
}

func newFileout(fileName string) (io.Writer, error) {
	logDir := filepath.Dir(fileName)
	if err := _create_file_dir(logDir); err != nil {
		return nil, err
	}
	if fw, err := os.OpenFile(fileName, (os.O_APPEND | os.O_RDWR | os.O_CREATE), os.ModePerm); err != nil {
		return nil, err
	} else {
		return bufio.NewWriter(fw), nil
	}
}
