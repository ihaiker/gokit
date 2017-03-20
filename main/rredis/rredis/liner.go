package rredis

import (
	"os"
	"github.com/peterh/liner"
	"strings"
	"fmt"
	"github.com/ihaiker/gokit/main/rredis/config"
)

func lineWordCompleter(cfg *config.Config) liner.WordCompleter {
	clients := make([]string, len(cfg.RemoteRedis))
	for idx, redis := range cfg.RemoteRedis {
		clients[idx] = redis.Name
	}
	cmdHead := make([]string, len(Commands))
	i := 0
	for name, _ := range Commands {
		cmdHead[i] = name
		i += 1
	}
	return func(line string, pos int) (head string, completions []string, tail string) {
		if strings.TrimSpace(line) == "" {
			return "", cmdHead, ""
		} else if strings.TrimSpace(line) == "use" || strings.HasPrefix("use", line) {
			return "use ", clients, ""
		}
		return line, []string{}, ""
	}

}

func Cmd(config *config.Config, serviceManager *ServiceManager) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)
	line.SetTabCompletionStyle(liner.TabCircular)
	line.SetWordCompleter(lineWordCompleter(config))

	if f, err := os.Open(config.HistoryPath); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	var err error
	cmd := ""
	LOOP:
	for {
		if cmd, err = line.Prompt("> "); err == nil {
			cmd = strings.Trim(cmd, " ")
			if cmd == "" {
				continue
			}
			args := strings.SplitN(cmd, " ", 2)
			execCommand, has := Commands[args[0]];
			if has {
				if len(args) == 1 {
					err = execCommand(serviceManager, config, "")
				} else {
					err = execCommand(serviceManager, config, args[1])
				}
			} else {
				err = defaultCommand(serviceManager, config, cmd);
			}
			if err != nil {
				if err == QUIT_ERROR {
					break LOOP
				} else {
					fmt.Println(err.Error())
				}
			}
			line.AppendHistory(cmd)
		} else if err == liner.ErrPromptAborted {
			fmt.Print("")
		} else {
			fmt.Print("Error reading line: ", err)
		}
	}

	if f, err := os.Create(config.HistoryPath); err != nil {
		fmt.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}