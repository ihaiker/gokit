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
)

var _loggers = make(map[string]*LoggerEntry)

const default_config = `
root:
    level: "info"
    appender: "console"
`

//获取某个logger的级别
func _level(logger string, cfg *config.Config) Level {
    //level
    if level, err := cfg.GetString(logger + ".level"); err != nil {
        panic(err)
    } else {
        return Level(strings.ToLower(level))
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

    {
        //root is default
        appender := _appender("root", cfg)
        log.SetOutput(appender)
        log.SetPrefix("[R] ")
        log.SetFlags(_LOG_FLAG)
    }

    config_logger := func(loggerName string) {
        logGroup := &LoggerEntry{}
        level := _level(loggerName, cfg)
        appender := _appender(loggerName, cfg)
        switch level {
        case _DEBUG:
            logGroup.debug_ = log.New(appender, "[D] ", _LOG_FLAG)
            fallthrough
        case _INFO:
            logGroup.info_ = log.New(appender, "[I] ", _LOG_FLAG)
            fallthrough
        case _WARN:
            logGroup.warn_ = log.New(appender, "[W] ", _LOG_FLAG)
            fallthrough
        case _ERROR:
            logGroup.error_ = log.New(appender, "[E] ", _LOG_FLAG)
        }
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

func init() {
    f := fileKit.New("./conf/logs.yaml")
    if f.Exist() {
        log.Println("use log config file ",f.GetPath())
        content, _ := f.ToString()
        if err := SetConfigWithContent(content); err != nil {
            log.Panic("set config :", err.Error())
        }
    } else {
        log.Println("the config file ",f.GetPath()," not found !")
        if err := SetConfigWithContent(""); err != nil {
            log.Panic("set config :", err.Error())
        }
    }
}

