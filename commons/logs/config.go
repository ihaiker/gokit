package logs

import (
    "os"
    "github.com/ihaiker/gokit/config"
    "github.com/ihaiker/gokit/files"
    "github.com/ihaiker/gokit/config/yaml"

    "errors"
    "io/ioutil"
    "io"
    "regexp"
    "log"
    "strings"
    "fmt"
    "github.com/spf13/pflag"
)

var _loggers map[string]*LoggerEntry

const default_config = `
root:
    level: "info"
    appender: "console"
`

var (
    colorOff   = "\033[0m"
    colorDebug = "\033[0;29m"
    colorInfo  = "\033[0;34m"
    colorWarn  = "\033[0;33m"
    colorError = "\033[0;31m"

    colorPath  = "\033[0;02m"
    colorClass = "\033[0;94m"
)

//获取某个logger的级别
func _level(logger string, cfg *config.Config) Level {
    //level
    if level, err := cfg.GetString(logger + ".level"); err != nil {
        panic(err)
    } else {
        return FromString(level)
    }
}

//输出地点
func _appender(logger string, cfg *config.Config) io.Writer {
    if appender, err := cfg.GetString(logger + ".appender"); err != nil {
        panic(err)
    } else {
        as := strings.SplitN(appender, ":", 2)
        switch as[0] {
        case "none":
            return ioutil.Discard
        case "console":
            return os.Stdout
        case "file":
            if len(as) == 1 {
                panic(logger + ".appender error. file:/path.{yyyy-MM-dd HH:mm:ss}")
            }
            path := as[1]
            match, _ := regexp.Match("\\{[yMdHmsS-]*\\}", []byte(path))
            if match {
                if out, err := NewDailyRollingFileOut(path); err != nil {
                    panic(err)
                } else {
                    return out
                }
            } else {
                if out, err := NewFileOut(path); err != nil {
                    panic(err)
                } else {
                    return out
                }
            }
        case "unix":
        case "sock":

        }
    }
    return ioutil.Discard
}
func _flag(logger string, cfg *config.Config) int {
    if flagStr, err := cfg.GetString(logger + ".flag"); err != nil || flagStr == "" {
        return _LOG_FLAG
    } else {
        var flag int
        for _, f := range strings.Split(flagStr, " ") {
            switch f {
            case "date":
                flag = flag | log.Ldate
            case "time":
                flag = flag | log.Ltime
                //case "longfile":
                //  flag = flag | log.Llongfile
                //case "shortfile":
                //    flag = flag | log.Lshortfile
            case "UTC":
                flag = flag | log.LUTC
            case "microseconds":
                flag = flag | log.Lmicroseconds
            }
        }
        return flag
    }
}

func SetConfig(configFile string) error {
    f := fileKit.New(configFile)
    if !f.Exist() {
        return errors.New("the config file not found! " + configFile)
    }
    content, err := f.ToString()
    if err != nil {
        return err
    }
    return SetConfigWithContent(content)
}

//设置日志配置项
func SetConfigWithContent(content string) (err error) {
    defer func() {
        if err == nil {
            if e := recover(); e != nil {
                err = e.(error)
            }
        }
    }()
    _loggers = make(map[string]*LoggerEntry)

    //init config
    var cfg *config.Config
    cfg, err = yaml.Config(default_config)
    if err != nil {
        err = errors.New("read default config error: " + err.Error())
        return
    }
    if content != "" {
        err = cfg.Load(content)
        if err != nil {
            err = errors.New("read config file error: " + err.Error())
            return
        }
    }

    config_logger := func(loggerName string) {
        logGroup := &LoggerEntry{}
        logGroup.level = _level(loggerName, cfg)
        appender := _appender(loggerName, cfg)
        flag := _flag(loggerName, cfg)

        if loggerName == "root" {
            log.SetOutput(appender)
            log.SetPrefix("root ")
            log.SetFlags(flag)
        }

        logGroup.debug_ = log.New(appender, fmt.Sprintf("%s[DEBUG] %s%s ", colorDebug, loggerName, colorOff), flag)
        logGroup.info_ = log.New(appender, fmt.Sprintf("%s[INFO] %s%s ", colorInfo, loggerName, colorOff), flag)
        logGroup.warn_ = log.New(appender, fmt.Sprintf("%s[WARN] %s%s ", colorWarn, loggerName, colorOff), flag)
        logGroup.error_ = log.New(appender, fmt.Sprintf("%s[ERROR] %s%s ", colorError, loggerName, colorOff), flag)
        _loggers[loggerName] = logGroup
    }
    config_logger("root")

    if defineLoggers, err := cfg.GetSlice("logger"); err == nil {
        for _, defineLogger := range defineLoggers {
            config_logger(defineLogger.(string))
        }
    }
    return
}

var debug = pflag.Bool("debug", false, "use debug module")
var log_config = pflag.String("logs-config", "", "the logs config file")

func init() {
    pflag.Parse()

    configFile := fileKit.New("./conf/logs.yaml")
    if *log_config != "" {
        configFile = fileKit.New(*log_config)
        if !configFile.Exist() {
            Fatal("日志文件不存在：", configFile.GetPath())
        }
    }

    if configFile.Exist() {
        content, _ := configFile.ToString()
        if err := SetConfigWithContent(content); err != nil {
            log.Panic("set config :", err.Error())
        }
    } else {
        if err := SetConfigWithContent(""); err != nil {
            log.Panic("set config :", err.Error())
        }
    }

    if *debug {
        SetAllLevel(DEBUG)
    }
}
