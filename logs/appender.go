package logs

import (
	"errors"
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
		return NewDailyRolling(file)

	} else if strings.HasPrefix(appender, "sock://") {

		return net.Dial("tcp", appender[7:])

	} else if strings.HasPrefix(appender, "unix://") {

		return net.Dial("unix", appender[6:])

	}
	return nil, errors.New("not support the appender:" + appender)
}
