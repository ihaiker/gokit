package logs

type Level int

const (
	DEBUG Level = 0
	INFO  Level = 1
	WARN  Level = 2
	ERROR Level = 3
	FATAL Level = 4
)

func (this *Level) Int() int {
	return int(*this)
}

func (this *Level) String() string {
	switch this.Int() {
	case 0:
		return "DEBUG"
	case 1:
		return "INFO"
	case 2:
		return "WARN"
	case 3:
		return "ERROR"
	default:
		return "FATAL"
	}
}

func (this *Level) Single() string {
	switch this.Int() {
	case 0:
		return "D"
	case 1:
		return "I"
	case 2:
		return "W"
	case 3:
		return "E"
	default:
		return "F"
	}
}

func (this *Level) PrintLevel(level Level) bool {
	return this.Int() <= level.Int()
}

func FromString(level string) Level {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return FATAL
	}
}
