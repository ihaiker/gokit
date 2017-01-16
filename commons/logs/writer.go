package logs

import (
	"github.com/bluele/gcache"
	"path/filepath"
	"github.com/ihaiker/gokit/files"
	"os"
	"time"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"
	"github.com/ihaiker/gokit/commons/time"
	"strings"
	"regexp"
)

type DailyRollingFile struct {
	cache   gcache.Cache
	pattern string
	layout  string
}

func (self *DailyRollingFile) Write(p []byte) (n int, err error) {
	key := time.Now().Format(self.layout)
	if fd, err := self.cache.Get(key); err != nil {
		return 0, err
	} else {
		w := fd.(io.Writer)
		return w.Write(p)
	}
}

func _create_file_dir(logDir string) error {
	if fileKit.NotExist(logDir) {
		if err := os.MkdirAll(logDir, os.ModeDir); err != nil {
			return errors.New("can mkdir " + logDir + " " + err.Error())
		}
	}
	return nil
}

func NewDailyRollingFileOut(dailyRollingFileName string) (io.Writer, error) {
	logDir := filepath.Dir(dailyRollingFileName)
	if err := _create_file_dir(logDir); err != nil {
		return nil, err
	}
	reg, _ := regexp.Compile("\\{[yMdHmsS-]*\\}")
	layout := string(reg.Find([]byte(dailyRollingFileName)))
	layout = layout[1:len(layout)-1]
	
	cb := gcache.New(10).LRU().Expiration(time.Hour * 24)
	cb.LoaderFunc(func(key interface{}) (interface{}, error) {
		return os.OpenFile(
			strings.Replace(dailyRollingFileName, "{" + layout + "}", key.(string), 1),
			(os.O_APPEND | os.O_RDWR | os.O_CREATE), os.ModePerm,
		)
	}).EvictedFunc(func(key, value interface{}) {
		fd := value.(*os.File)
		fd.Close()
	})
	
	return &DailyRollingFile{
		cache:cb.Build(),
		pattern:layout,
		layout:timeKit.GoLayout(layout),
	}, nil
}

func NewFileOut(fileName string) (io.Writer, error) {
	logDir := filepath.Dir(fileName)
	if err := _create_file_dir(logDir); err != nil {
		return nil, err
	}
	return os.OpenFile(fileName, (os.O_APPEND | os.O_RDWR | os.O_CREATE), os.ModePerm)
}