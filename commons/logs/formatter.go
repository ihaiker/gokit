package logs

import (
	"github.com/Sirupsen/logrus"
	"bytes"
	"time"
	"strings"
	"github.com/ihaiker/gokit/commons/time"
	"gopkg.in/square/go-jose.v1/json"
	"fmt"
	"sort"
)

type FormatField uint8

const (
	LEVEL FormatField = iota
	DATETIME
	DATA_FIELDS
	MESSAGE
	LF
)

const logger_path = "/github.com/Sirupsen/logrus/logger.go"
const entry_path = "/github.com/Sirupsen/logrus/entry.go"

var found = []slice{
	slice{key:"%p", field:LEVEL, replace:"%s"},
	slice{key:"%d", field:DATETIME, replace:"%s"},
	slice{key:"%f", field:DATA_FIELDS, replace:"%s"},
	slice{key:"%m", field:MESSAGE, replace:"%s"},
	slice{key:"%n", field:LF, replace:"\n"},
}

type slice struct {
	key     string
	index   int
	field   FormatField
	replace string
}
type slices []slice

func (self slices) Len() int {
	return len(self)
}
func (self slices) Less(i, j int) bool {
	s1 := self[i]
	s2 := self[j]
	return s1.index < s2.index
}
func (self slices) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type PatternFormatter struct {
	pattern    string
	fields     slices
	dateLayout string
}

func (self *PatternFormatter) datetime(time time.Time) string {
	return time.Format(self.dateLayout)
}

func (self *PatternFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	values := []interface{}{}
	buf := bytes.NewBuffer([]byte{})
	for _, f := range self.fields {
		switch f.field {
		case LEVEL:
			values = append(values, entry.Level.String())
		case MESSAGE:
			values = append(values, entry.Message)
		case DATETIME:
			values = append(values, self.datetime(entry.Time))
		case DATA_FIELDS:
			bs, _ := json.Marshal(entry.Data)
			values = append(values, string(bs))
		}
	}
	buf.WriteString(fmt.Sprintf(self.pattern, values...))
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

func NewFormatter(pattern string) *PatternFormatter {
	fields := slices{}
	for _, f := range found {
		for idx := strings.Index(pattern, f.key); idx != -1; idx = strings.Index(pattern, f.key) {
			f.index = idx
			fields = append(fields, f)
			pattern = strings.Replace(pattern, f.key, f.replace, 1)
		}
	}
	sort.Sort(fields)
	return &PatternFormatter{
		pattern:pattern,
		fields:fields,
		dateLayout:timeKit.GoLayout("yyyy-MM-dd HH:mm:ss.SSS"),
	}
}