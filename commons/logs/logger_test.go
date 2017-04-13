package logs

import (
    "testing"
    "log"
    "path/filepath"
)

func init() {
    p, _ := filepath.Abs(".")
    if err := SetConfig(p + "/logs.yaml"); err != nil {
        log.Print(err)
    }
}

func TestSetConfig(t *testing.T) {
    log.Println("aaa", "bb")

    Debug("DEBUG")
    Info("INFO")
    Warn("WARN")
    Error("ERROR")

    Debugf("%s,,","debugf")
    Infof("%s,,","infof")
    Warnf("%s,,","warnf")
    Errorf("%s,,","errorf")
}

func TestLogger(t *testing.T) {
    logger := Logger("console")
    logger.Debug("console debug")
    logger.Info("console info")
    logger.Warn("console warn")
    logger.Error("console error")
}


func TestLoggerF(t *testing.T) {
    logger := Logger("console")
    logger.Debugf("console debug %s.%d","f",1)
    logger.Infof("console info %s.%d","f",2)
    logger.Warnf("console warn %s.%d","f",3)
    logger.Errorf("console error %s.%d","f",4)
}