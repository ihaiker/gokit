package logs

import (
	"fmt"
	"github.com/ihaiker/gokit/config"
	"github.com/jinzhu/configor"
	"io"
	"os"
	"path/filepath"
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
		With(&configor.Config{
			ENVPrefix: strings.ToUpper(appName) + "_LOG",
			Silent:    true,
		}).Marshal(loggerConfiger)
	return
}

func initLogger(name string, itemConfig *loggerConfigItem) error {
	logger := GetLogger(name)

	logger.SetLevel(FromString(itemConfig.Level))
	logger.SetPattern(itemConfig.Pattern)
	if appenderWriter, err := appender(itemConfig.Appender); err != nil {
		return err
	} else {
		logger.SetOut(appenderWriter)
	}

	return nil
}

func Open(appName string) {
	cfg, err := loadConfig(appName)
	if err != nil {
		fmt.Println("load config error: ", err)
	}

	if cfg.Root != nil {
		if err = initLogger("root", cfg.Root); err != nil {
			fmt.Println("init root logger error:", err)
		}
	}
	if cfg.Root == nil || err != nil {
		_ = initLogger("root", &loggerConfigItem{
			Level:    "info",
			Appender: "stdout",
			Pattern:  DEFAULT_PATTERN,
		})
	}

	for name, logger := range cfg.Loggers {
		if err := initLogger(name, logger); err != nil {
			fmt.Printf("init %s logger error: %s", name, err)
		}
	}
}

func CloseAll() {
	for _, logger := range loggers {
		out := logger.Out()
		if out == os.Stdout {
			continue
		}
		if closer, match := out.(io.Closer); match {
			_ = closer.Close()
		}
	}
}

func init() {
	Open(filepath.Base(os.Args[0]))
}
