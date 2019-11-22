package logs

import (
	"github.com/ihaiker/gokit/config"
	"github.com/jinzhu/configor"
	"os"
	"path"
	"strings"
)

type loggerConfigItem struct {
	Level    string `json:"level" yaml:"level"`
	Appender string `json:"appender" yaml:"appender"`
	Pattern  string `json:"pattern" yaml:"pattern"`
}

type loggerConfig struct {
	Root    *loggerConfigItem            `json:"root" yaml:"root"`
	Loggers map[string]*loggerConfigItem `json:"loggers" yaml:"loggers"`
}

func loadConfig(appName string) (loggerConfiger *loggerConfig, err error) {
	loggerConfiger = &loggerConfig{
		Root: &loggerConfigItem{
			Level:    "info",
			Appender: "stdout",
			Pattern:  DEFAULT_PATTERN,
		},
	}
	err = config.NewConfigRegister(appName, "logs").
		With(&configor.Config{ENVPrefix: strings.ToUpper(appName) + "_LOG"}).Marshal(loggerConfiger)

	//忽略配置没有找到错误
	if err == config.ErrConfigNotFound {
		err = nil
	}

	if err != nil {
		return
	}
	return
}

func init() {
	appName := path.Base(os.Args[0])
	if strings.HasSuffix(appName, ".exe") { //windows
		appName = appName[:len(appName)-4]
	}
	if err := InitLoggers(appName); err != nil {
		panic(err)
	}
}

func initLogger(name string, itemConfig *loggerConfigItem) error {
	logger, has := Log(name)
	if ! has {
		logger = createLogger(name)
	}
	logger.(ConfigLogger).SetLevel(FromString(itemConfig.Level))
	logger.(ConfigLogger).SetPattern(itemConfig.Pattern)
	if appenderWriter, err := appender(itemConfig.Appender); err != nil {
		return err
	} else {
		logger.(ConfigLogger).SetOut(appenderWriter)
	}

	loggers[name] = logger
	return nil
}

func InitLoggers(appName string) (error) {
	cfg, err := loadConfig(appName)
	if err != nil {
		return err
	}

	if err := initLogger("root", cfg.Root); err != nil {
		return err
	}

	for name, logger := range cfg.Loggers {
		if err := initLogger(name, logger); err != nil {
			return err
		}
	}
	return nil
}
