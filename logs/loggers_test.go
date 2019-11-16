package logs

import (
	"testing"
)

func TestLogger(t *testing.T) {
	Root().Debug("root debug")
	Root().Info("root info")
	Root().Warn("root warnings")
	Root().Error("root ERROR")
	Root().Fatal("root fatal")
	Root().Infof("message: %s", "root warnings")

	test := GetLogger("ctl")
	test.Debug("debug =======")
	test.Info("info ========")
	test.Warn("warn =========")

	Debug("=========== root debug")
	Info("=========== root info")
}
