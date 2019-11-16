package logs

import (
	"errors"
	fileKit "github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/logs/appenders"
	"io"
	"net"
	"os"
	"strings"
)

func appender(appender string) (io.Writer, error) {
	if appender == "stdout" || appender == "console" {

		return os.Stdout, nil

	} else if strings.HasPrefix(appender, "file://") {
		file := appender[6:]
		if _, match := appenders.MatchDailyRollingFile(file); match {
			return appenders.NewDailyRollingFileOut(file)
		}
		return fileKit.New(file).GetWriter(true)

	} else if strings.HasPrefix(appender, "sock://") {

		return net.Dial("tcp", appender[7:])

	} else if strings.HasPrefix(appender, "unix://") {

		return net.Dial("unix", appender[6:])

	}
	return nil, errors.New("not support the appender:" + appender)
}
