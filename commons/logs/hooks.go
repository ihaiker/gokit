package logs

import (
	"io"
	"github.com/Sirupsen/logrus"
)

type fsHook struct {
	out       io.Writer
	formatter logrus.Formatter
	levels    []logrus.Level
}

func (self *fsHook) Fire(entry *logrus.Entry) error {
	if bs, err := self.formatter.Format(entry); err != nil {
		return err
	} else {
		_, err = self.out.Write(bs)
		return err
	}
}

func (self *fsHook) Levels() []logrus.Level {
	return self.levels
}