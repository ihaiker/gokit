package logs

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const DEFAULT_PATTERN = "[%L]%d{2006-01-02 15:04:05} [%a] %f:%l %F %m"

type content func(entry *entry) []byte

type pattern struct {
	layout   string
	contents []content
}

func (p *pattern) write(out io.Writer, entry *entry) {
	bufter := new(bytes.Buffer)
	for _, method := range p.contents {
		_, _ = bufter.Write(method(entry))
	}
	bufter.WriteRune('\n')
	_, _ = out.Write(bufter.Bytes())
}

func (p *pattern) String() string {
	return p.layout
}

func stringToken(obj ...interface{}) content {
	return func(entry *entry) []byte {
		return obj[0].([]byte)
	}
}

func nameToken(obj ...interface{}) content {
	return func(entry *entry) []byte {
		return []byte(entry.name)
	}
}

func levelToken(obj ...interface{}) content {
	param := ""
	if len(obj) > 0 {
		param = fmt.Sprint(obj...)
	}
	return func(entry *entry) []byte {
		if param == "1" {
			return []byte(entry.level.Single())
		} else {
			return []byte(entry.level.String())
		}
	}
}

func timeToken(obj ...interface{}) content {
	layout := "yyyy-MM-dd HH:mm:ss.SSS"
	if len(obj) == 1 {
		layout = obj[0].(string)
	}
	return func(entry *entry) []byte {
		return []byte(entry.time.Format(layout))
	}
}
func fileToken(obj ...interface{}) content {
	return func(entry *entry) []byte {
		return []byte(entry.file)
	}
}

func lineToken(obj ...interface{}) content {
	return func(entry *entry) []byte {
		return []byte(strconv.Itoa(entry.line))
	}
}

func funcToken(obj ...interface{}) content {
	return func(entry *entry) []byte {
		return []byte(entry.fun)
	}
}

func messageToken(obj ...interface{}) content {
	return func(entry *entry) []byte {
		return []byte(entry.message)
	}
}

//[%L] %d{yyyy-MM-dd HH:mm:ss.SSS} %f:%l %F %m
//%L{n} 级别编号，
func newPattern(layout string) (*pattern) {
	idxA := strings.Index(layout, "%a")
	idxL := strings.Index(layout, "%L")
	idxd := strings.Index(layout, "%d")
	idxf := strings.Index(layout, "%f")
	idxl := strings.Index(layout, "%l")
	idxF := strings.Index(layout, "%F")
	idxm := strings.Index(layout, "%m")
	idxMap := map[int](func(...interface{}) content){
		idxA: nameToken,
		idxL: levelToken, idxd: timeToken, idxf: fileToken,
		idxl: lineToken, idxF: funcToken, idxm: messageToken,
	}
	pattern := &pattern{contents: []content{}, layout: layout}

	bs := []byte(layout)

	var st []byte
	for i := 0; i < len(bs); {
		if fun, has := idxMap[i]; has {
			if st != nil {
				pattern.contents = append(pattern.contents, stringToken(st))
				st = nil
			}
			i++
			if i >= len(bs) {
				break
			}
			paramStartIdx := i + 1
			if len(bs) > paramStartIdx && layout[paramStartIdx:paramStartIdx+1] == "{" {
				if paramEndIdx := strings.Index(layout[paramStartIdx+1:], "}"); paramEndIdx == -1 {
					fmt.Println("log pattern error, use default:", DEFAULT_PATTERN)
					return newPattern(DEFAULT_PATTERN)
				} else {
					i = paramStartIdx + paramEndIdx + 1
					param := layout[paramStartIdx+1 : paramStartIdx+paramEndIdx+1]
					pattern.contents = append(pattern.contents, fun(param))
				}
			} else {
				pattern.contents = append(pattern.contents, fun())
			}
		} else if st == nil {
			st = []byte{bs[i]}
		} else {
			st = append(st, bs[i])
		}
		i++
	}
	if st != nil {
		pattern.contents = append(pattern.contents, stringToken(st))
	}
	return pattern
}
