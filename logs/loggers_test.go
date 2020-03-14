package logs

import (
	"io"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	defer CloseAll()

	Root().Debug("root debug")
	Root().Info("root info")
	Root().Warn("root warnings")
	Root().Error("root ERROR")
	Root().Fatal("root fatal")
	Root().Infof("message: %s", "root warnings")

}

func TestDef(t *testing.T) {
	Debug("=========== root debug")
	Info("=========== root info")
}

func TestGetLogger(t *testing.T) {
	defer CloseAll()
	SetDebugMode(true)
	test := GetLogger("ctl")
	test.Debug("debug =======")
	test.Info("info ========")
	test.Warn("warn =========")
}

func TestDailyRollingFileOutqq(t *testing.T) {
	filename := "/tmp/m{20060102}/stdout.{200601021504}.log"
	w, err := NewDailyRolling(filename)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond)
		_, _ = w.Write([]byte(time.Now().Format(time.RFC3339) + "\n"))
	}
	t.Log(w.(io.Closer).Close())
}
