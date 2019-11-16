package logs

import (
	"os"
	"testing"
	"time"
)

func TestNewPattern(t *testing.T) {
	p := newPattern("([%L] %d{2006-01-02 15:04:05} %F:%l %m)")
	out := os.Stdout
	entry := &entry{
		level: INFO, //%L
		time: time.Now(), //%d
		file:    "github.com/ihaiker/gokit/gokit.go",     //%f
		line:    10,             //%l
		fun:     "gokit.Version", //%F
		message: "商贸内容",         //%m
	}
	p.write(out, entry)
}
