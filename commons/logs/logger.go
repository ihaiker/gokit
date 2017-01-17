package logs

import (
	"github.com/Sirupsen/logrus"
	"os"
	"github.com/ihaiker/gokit/config"
	"github.com/ihaiker/gokit/files"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io/ioutil"
	"io"
	"github.com/ihaiker/gokit/config/yaml"
	"regexp"
)

const default_config = `
root:
    level: "debug"
    appender: "console"
    formatter: "pattern"
    layout: "[%p] %d %f %m"
`

func _level(logger string, cfg *config.Config) logrus.Level {
	//level
	if level, err := cfg.GetString(logger + ".level"); err != nil {
		panic(err)
	} else {
		if levelObj, err := logrus.ParseLevel(level); err != nil {
			panic(err)
		} else {
			return levelObj
		}
	}
}
func _appender(logger string, cfg *config.Config) io.Writer {
	if appender, err := cfg.GetString(logger + ".appender"); err != nil {
		panic(err)
	} else {
		switch appender {
		case "none":
			return ioutil.Discard
		case "console":
			return os.Stdout
		case "file":
			path, err := cfg.GetString(logger + ".path")
			if err != nil {
				panic(err)
			}
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
		}
	}
	return ioutil.Discard
}
func _formatter(logger string, cfg *config.Config) logrus.Formatter {
	if formatter, err := cfg.GetString(logger + ".formatter"); err != nil {
		panic(err)
	} else {
		switch formatter {
		case "text":
			return new(logrus.TextFormatter)
		case "json":
			return new(logrus.JSONFormatter)
		case "pattern":
			if layout, err := cfg.GetString(logger + ".layout"); err != nil {
				panic(err)
			} else {
				return NewFormatter(layout)
			}
		default:
			panic(errors.New("not found formatter" + formatter))
		}
	}
}

func parseHook(logger string, cfg *config.Config) logrus.Hook {
	level := _level(logger, cfg)
	levels := []logrus.Level{}
	for _, l := range logrus.AllLevels {
		if l <= level {
			levels = append(levels, l)
		}
	}
	
	hook := &fsHook{
		out:_appender(logger, cfg), levels:levels,
		formatter:_formatter(logger, cfg),
	}
	return hook
}

func SetDefault() error {
	return SetConfigWithContent("")
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
	
	logrus.SetLevel(_level("root", cfg))
	logrus.SetOutput(_appender("root", cfg))
	logrus.SetFormatter(_formatter("root", cfg))
	
	//清楚之前hook
	for level,_ := range logrus.StandardLogger().Hooks {
		delete(logrus.StandardLogger().Hooks,level)
	}
	if loggers, err := cfg.GetSlice("logger"); err != nil {
		if err != config.NOT_FOUND {
			return err
		}
	} else {
		for _, log := range loggers {
			logrus.AddHook(parseHook(log.(string), cfg))
		}
	}
	return
}

func init() {
	SetDefault()
}