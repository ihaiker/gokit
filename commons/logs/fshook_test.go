package logs

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	log,err := New("logs.yaml")
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}
	log.Debugf("a - %d", 12)
}