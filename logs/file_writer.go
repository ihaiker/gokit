package logs

import (
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/files"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var compile = regexp.MustCompile(`\{[0-9:-]*\}`)

//每日回滚日志
type dailyRollingFile struct {
	layout string
	c      chan []byte

	current string
	out     io.Writer
	timer   *time.Timer
	gw      *sync.WaitGroup
}

func (self *dailyRollingFile) start() {
	defer func() {
		for p := range self.c {
			_, _ = self.out.Write(p)
		}
		self.gw.Done()
	}()

	for {
		select {
		case <-self.timer.C:
			self.timer.Reset(time.Second)
			if self.current != self.rollingFilename(time.Now()) {
				if err := self.rollingFile(); err != nil {
					fmt.Println("rollingFile file ", self.layout, " error : ", err)
				}
			}
		case line, has := <-self.c:
			if !has {
				return
			}
			_, _ = self.out.Write(line)
		}
	}
}

func (self *dailyRollingFile) Write(p []byte) (n int, err error) {
	defer errors.Catch(func(e error) {
		err = e
	})
	self.c <- p
	return len(p), nil
}

func (self *dailyRollingFile) Close() error {
	errors.Exec(func() { close(self.c) })
	self.gw.Wait()
	self.closePre()
	return nil
}

func (self *dailyRollingFile) closePre() {
	if self.out == nil {
		return
	}
	if self.out == os.Stdout {

	} else if out, match := self.out.(io.Closer); match {
		_ = out.Close()
	}
}

func (self *dailyRollingFile) rollingFilename(t time.Time) string {
	file := self.layout
	for {
		if find := compile.FindString(file); find == "" {
			break
		} else {
			format := find[1 : len(find)-1]
			file = strings.Replace(file, find, t.Format(format), 1)
		}
	}
	return file
}

func (self *dailyRollingFile) rollingFile() (err error) {
	self.closePre()
	now := time.Now()
	self.current = self.rollingFilename(now)

	if err = os.MkdirAll(filepath.Dir(self.current), os.ModePerm); err != nil {
		return
	}
	self.out, err = os.OpenFile(self.current, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		self.out = os.Stdout
	}
	return err
}

func NewDailyRolling(pattern string) (io.Writer, error) {
	if ! compile.MatchString(pattern) {
		return files.New(pattern).GetWriter(true)
	}
	daily := &dailyRollingFile{
		layout: pattern, gw: new(sync.WaitGroup), c: make(chan []byte, 100),
		timer: time.NewTimer(time.Second),
	}
	if err := daily.rollingFile(); err != nil {
		return nil, err
	}
	daily.gw.Add(1)
	go daily.start()
	return daily, nil
}
