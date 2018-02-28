package logs

import (
    "testing"
    "log"
    "path/filepath"
    "fmt"
    "github.com/fatih/color"
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
    logger.Debug("console","debug")
    logger.Info("console","info")
    logger.Warn("console","warn")
    logger.Error("console","error")
}


func TestLoggerF(t *testing.T) {
    logger := Logger("console")
    logger.Debugf("console debug %s.%d","f",1)
    logger.Infof("console info %s.%d","f",2)
    logger.Warnf("console warn %s.%d","f",3)
    logger.Errorf("console error %s.%d","f",4)
}

func TestWarn(t *testing.T) {
    yellow := color.New(color.FgYellow).SprintFunc()
    red := color.New(color.FgRed).SprintFunc()
    fmt.Printf("This is a %s and this is %s.\n", yellow("warning"), red("error"))

    info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
    fmt.Printf("This %s rocks!\n", info("package"))

    // Use helper functions
    fmt.Println("This", color.RedString("warning"), "should be not neglected.")
    fmt.Printf("%v %v\n", color.GreenString("Info:"), "an important message.")

    // Windows supported too! Just don't forget to change the output to color.Output
    fmt.Fprintf(color.Output, "Windows support: %s", color.GreenString("PASS"))
}