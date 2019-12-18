package appenders

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/ihaiker/gokit/files"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type ClosedBufIOWriter struct {
	writer    io.Writer
	bufWriter io.Writer
}

func NewClosedBufIOWriter(writer io.Writer) io.WriteCloser {
	return &ClosedBufIOWriter{
		writer:    writer,
		bufWriter: bufio.NewWriter(writer),
	}
}
func (self *ClosedBufIOWriter) Write(p []byte) (n int, err error) {
	n, err = self.bufWriter.Write(p)
	return
}

func (self *ClosedBufIOWriter) Close() error {
	if closed, match := self.writer.(io.Closer); match {
		return closed.Close()
	}
	return nil
}

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
		if cl, match := out.(io.Closer); match {
			_ = cl.Close()
		}
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
		return NewClosedBufIOWriter(fw), nil
	}
}
