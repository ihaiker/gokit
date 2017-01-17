package logs

import (
	"testing"
)

func TestDefault(t *testing.T) {
	err := SetDefault()
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}
	Debugf("a - %d", 12)
}

func TestNewLogger(t *testing.T) {
	err := SetConfig("logs.yaml")
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}
	Debugf("a - %d", 12)
}